package app

import (
	"Notes/config"
	_ "Notes/docs"
	"Notes/internal/api/http/handler"
	"Notes/internal/api/http/middleware"
	"Notes/internal/repository"
	"Notes/internal/service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Collection struct {
	Auth     *handler.AuthHandler
	User     *handler.UserHandler
	Folder   *handler.FolderHandler
	Notebook *handler.NotebookHandler
	Note     *handler.NoteHandler
}

type Dependencies struct {
	SQL              *sql.DB
	Handlers         Collection
	AuthMiddleware   gin.HandlerFunc
	LoggerMiddleware gin.HandlerFunc
}

func Run() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	deps, err := setupDependencies(cfg)
	if err != nil {
		log.Fatalf("Dependencies setup failed: %v", err)
	}
	defer deps.SQL.Close()

	router := setupRouter(deps.Handlers, deps.AuthMiddleware, deps.LoggerMiddleware)
	srv := startHTTPServer(router, cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
	log.Println("Server exited properly")
}

func connectAndMigrate(cfg config.Database) (*gorm.DB, *sql.DB) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения *sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации мигратора: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	log.Println("Миграции успешно применены")

	for i := 0; i < 10; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		log.Println("БД не отвечает, повтор через 1 секунду...")
		time.Sleep(time.Second)
	}

	return db, sqlDB
}

func setupDependencies(cfg *config.Config) (*Dependencies, error) {
	gormDb, sqlDb := connectAndMigrate(cfg.Database)

	postgresRepo := repository.NewPostgresRepository(gormDb)
	jwtService := service.NewConcreteJwtService(cfg)
	hashService := service.NewConcreteHashService()
	authService := service.NewConcreteAuthService(postgresRepo, jwtService)
	folderService := service.NewConcreteFolderService(postgresRepo)
	notebookService := service.NewConcreteNotebookService(postgresRepo)
	noteService := service.NewConcreteNoteService(postgresRepo)
	userService := service.NewConcreteUserService(postgresRepo, hashService)

	return &Dependencies{
		SQL: sqlDb,
		Handlers: Collection{
			Auth:     handler.NewAuthHandler(authService),
			User:     handler.NewUserHandler(userService),
			Folder:   handler.NewFolderHandler(folderService),
			Notebook: handler.NewNotebookHandler(notebookService),
			Note:     handler.NewNoteHandler(noteService),
		},
		AuthMiddleware:   middleware.AuthMiddleware(authService),
		LoggerMiddleware: middleware.RequestLogger(),
	}, nil
}

func startHTTPServer(handler http.Handler, port int) *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Printf("Server started on port %d", port)
	return srv
}

func setupRouter(h Collection, authMiddleware gin.HandlerFunc, loggerMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(authMiddleware)
	r.Use(loggerMiddleware)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/user", h.User.GetUser)
		protected.PUT("/user", h.User.UpdateUser)
		protected.DELETE("/user", h.User.DeleteUser)

		protected.POST("/folder", h.Folder.CreateFolder)
		protected.PUT("/folder/:id", h.Folder.UpdateFolder)
		protected.DELETE("/folder/:id", h.Folder.DeleteFolder)

		protected.GET("/notebook", h.Notebook.GetNotebook)

		protected.POST("/notes", h.Note.CreateNote)
		protected.PUT("/notes/:id", h.Note.UpdateNote)
		protected.DELETE("/notes/:id", h.Note.DeleteNote)
		protected.GET("/notes/favorites", h.Note.GetFavoriteNotes)
		protected.GET("/notes/search", h.Note.FindNotes)
		protected.PUT("/notes/:id/move", h.Note.MoveNote)
		protected.PUT("/notes/:id/favorites", h.Note.AddToFavorites)
		protected.DELETE("/notes/:id/favorites", h.Note.DeleteFromFavorites)
	}

	r.POST("/api/auth/login", h.Auth.Login)
	r.POST("/api/user", h.User.CreateUser)

	return r
}
