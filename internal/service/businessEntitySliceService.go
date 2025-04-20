package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
	"math/rand"
	"time"
)

const entitiesCount = 3
const milliSecondsToSleep = 400

type BusinessEntitySliceService struct {
	repository repository.AbstractRepository
}

func NewBusinessEntitySliceService(repo repository.AbstractRepository) *BusinessEntitySliceService {
	return &BusinessEntitySliceService{
		repository: repo,
	}
}

func (s *BusinessEntitySliceService) GenerateAndSoreEntities(iterations int) {
	for i := 0; i < iterations; i++ {
		entity := s.generateRandomEntity()
		s.repository.SaveEntity(entity)
		time.Sleep(milliSecondsToSleep * time.Millisecond)
	}
}

func (s *BusinessEntitySliceService) generateRandomEntity() model.BusinessEntity {
	switch rand.Intn(entitiesCount) {
	case 0:
		newUser, _ := model.NewUser("Name", "Surname", "Login4321", "Password1234$")
		return newUser
	case 1:
		newFolder, _ := model.NewFolder("Title", 1)
		return newFolder
	default:
		newNote, _ := model.NewNote("Title", "Sample text", 1, nil)
		return newNote
	}
}
