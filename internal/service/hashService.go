package service

//go:generate mockgen -source=hashService.go -destination=mock/hashService.go -package=mock

import (
	"Notes/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AbstractHashService interface {
	GetHash(stringToHash string) (string, *model.ApplicationError)
}

type ConcreteHashService struct {
}

func NewConcreteHashService() AbstractHashService {
	return &ConcreteHashService{}
}

func (h *ConcreteHashService) GetHash(stringToHash string) (string, *model.ApplicationError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(stringToHash), bcrypt.DefaultCost)
	if err != nil {
		return "", model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при создании хэша", err)
	}
	return string(hashedPassword), nil
}
