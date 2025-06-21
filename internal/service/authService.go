package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
)

//go:generate mockgen -source=authService.go -destination=mock/authService.go -package=mock

type AbstractAuthService interface {
	AuthUser(login, password string) (string, *model.ApplicationError)
	ValidateToken(token string) (*model.Claims, *model.ApplicationError)
}

type ConcreteAuthService struct {
	repo repository.AbstractRepository
	jwt  AbstractJwtService
}

func NewConcreteAuthService(repository repository.AbstractRepository, jwtService AbstractJwtService) AbstractAuthService {
	return &ConcreteAuthService{
		repo: repository,
		jwt:  jwtService,
	}
}

func (a *ConcreteAuthService) AuthUser(login, password string) (string, *model.ApplicationError) {
	user, err := a.repo.GetUser(login, password)

	if err != nil {
		return "", err
	}

	return a.jwt.GetToken(user.Id)
}

func (a *ConcreteAuthService) ValidateToken(token string) (*model.Claims, *model.ApplicationError) {
	claims, err := a.jwt.ParseToken(token)

	if err != nil {
		return nil, err
	}

	_, errGetUser := a.repo.GetUserById(claims.UserId)

	if errGetUser != nil {
		return nil, errGetUser
	}
	return claims, nil
}
