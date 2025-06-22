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
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log.Println(cfg)

	gormDb, sqlDb := ConnectAndMigrate(cfg.Database)

	defer sqlDb.Close()

	postgresRepo := repository.NewPostgresRepository(gormDb)
	jwtService := service.NewConcreteJwtService(cfg)
	hashService := service.NewConcreteHashService()
	authService := service.NewConcreteAuthService(postgresRepo, jwtService)
	folderService := service.NewConcreteFolderService(postgresRepo)
	notebookService := service.NewConcreteNotebookService(postgresRepo)
	noteService := service.NewConcreteNoteService(postgresRepo)
	userService := service.NewConcreteUserService(postgresRepo, hashService)

	authHandler := handler.NewAuthHandler(authService)
	folderHandler := handler.NewFolderHandler(folderService)
	notebookHandler := handler.NewNotebookHandler(notebookService)
	noteHandler := handler.NewNoteHandler(noteService)
	userHandler := handler.NewUserHandler(userService)

	//logger.InitLogger(postgresRepo)
	ctx, cancelLogger := context.WithCancel(context.Background())
	defer cancelLogger()
	//logger.LogEntitiesAsync(postgresRepo, ctx)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.AuthMiddleware(authService)

	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/user", userHandler.GetUser)
		protected.PUT("/user", userHandler.UpdateUser)
		protected.DELETE("/user", userHandler.DeleteUser)

		protected.POST("/folder", folderHandler.CreateFolder)
		protected.PUT("/folder{id}", folderHandler.UpdateFolder)
		protected.DELETE("/folder{id}", folderHandler.DeleteFolder)

		protected.GET("/notebook", notebookHandler.GetNotebook)

		protected.POST("/notes", noteHandler.CreateNote)
		protected.PUT("/notes{id}", noteHandler.UpdateNote)
		protected.DELETE("/notes{id}", noteHandler.DeleteNote)
		protected.GET("/notes/favorites", noteHandler.GetFavoriteNotes)
		protected.GET("/notes/search", noteHandler.FindNotes)
		protected.PUT("/notes/{id}/move", noteHandler.MoveNote)
		protected.PUT("/notes/{id}/favorites", noteHandler.AddToFavorites)
		protected.DELETE("/notes/{id}/favorites", noteHandler.DeleteFromFavorites)
	}

	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/user", userHandler.CreateUser)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %s", err)
		}
	}()

	log.Printf("Server started on port %d", cfg.Server.Port)

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

func ConnectAndMigrate(cfg config.Database) (*gorm.DB, *sql.DB) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{LogLevel: logger.Info},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
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
		log.Fatalf("❌ Ошибка инициализации мигратора: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("❌ Ошибка применения миграций: %v", err)
	}

	log.Println("✅ Миграции успешно применены")

	for i := 0; i < 10; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		log.Println("БД не отвечает, повтор через 1 секунду...")
		time.Sleep(time.Second)
	}

	return db, sqlDB
}
