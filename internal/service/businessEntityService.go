package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
	"math/rand"
	"time"
)

const entitiesCount = 3
const milliSecondsToSleep = 200

type BusinessEntityService struct {
	repository repository.AbstractRepository
}

func NewBusinessEntityService(repo repository.AbstractRepository) *BusinessEntityService {
	return &BusinessEntityService{
		repository: repo,
	}
}

func (s *BusinessEntityService) GenerateEntitiesAsync(ch chan<- model.BusinessEntity, stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(milliSecondsToSleep * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ch <- s.generateRandomEntity()
			case <-stop:
				return
			}
		}
	}()
}

func (s *BusinessEntityService) SaveEntitiesAsync(ch <-chan model.BusinessEntity, stop <-chan struct{}) {
	go func() {
		for {
			select {
			case entity := <-ch:
				s.repository.SaveEntity(entity)
			case <-stop:
				return
			}
		}
	}()
}

func (s *BusinessEntityService) generateRandomEntity() model.BusinessEntity {
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
