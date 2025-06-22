package app

import (
	"Notes/config"
	_ "Notes/docs"
	"Notes/internal/api/http/handler"
	"Notes/internal/api/http/middleware"
	"Notes/internal/logger"
	"Notes/internal/repository"
	"Notes/internal/service"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	repo := repository.NewPlainRepository()
	jwtService := service.NewConcreteJwtService(cfg)
	hashService := service.NewConcreteHashService()
	authService := service.NewConcreteAuthService(repo, jwtService)
	folderService := service.NewConcreteFolderService(repo)
	notebookService := service.NewConcreteNotebookService(repo)
	noteService := service.NewConcreteNoteService(repo)
	userService := service.NewConcreteUserService(repo, hashService)

	authHandler := handler.NewAuthHandler(authService)
	folderHandler := handler.NewFolderHandler(folderService)
	notebookHandler := handler.NewNotebookHandler(notebookService)
	noteHandler := handler.NewNoteHandler(noteService)
	userHandler := handler.NewUserHandler(userService)

	logger.InitLogger(repo)
	ctx, cancelLogger := context.WithCancel(context.Background())
	defer cancelLogger()
	logger.LogEntitiesAsync(repo, ctx)

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
