package main

import (
	"Notes/internal/logger"
	"Notes/internal/model"
	"Notes/internal/repository"
	"Notes/internal/service"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	stop := make(chan struct{})
	entitiesChan := make(chan model.BusinessEntity)

	repo := repository.NewPlainRepository()
	businessService := service.NewBusinessEntityService(repo)
	logger.InitLogger(repo)

	businessService.GenerateEntitiesAsync(entitiesChan, stop)
	businessService.SaveEntitiesAsync(entitiesChan, stop)
	logger.LogEntitiesAsync(repo, stop)

	fmt.Println("Business entities generation and saving is running. Press Ctrl+C to stop")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Stopping service.")

	close(stop)
	time.Sleep(2 * time.Second)
}
