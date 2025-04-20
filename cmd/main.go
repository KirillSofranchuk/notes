package main

import (
	"Notes/internal/repository"
	"Notes/internal/service"
)

func main() {
	repo := repository.NewSliceRepository()
	businessEntityService := service.NewBusinessEntitySliceService(repo)
	businessEntityService.GenerateAndSoreEntities(22)
}
