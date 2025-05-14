package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
)

type AbstractAuthService interface {
	AuthUser(login, password string) (string, *model.ApplicationError)
	ValidateToken(token string) (*model.Claims, *model.ApplicationError)
}

type ConcreteAuthService struct {
	repo repository.AbstractRepository
	jwt  JwtService
}

func NewConcreteAuthService(repository repository.AbstractRepository, jwtService JwtService) AbstractAuthService {
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
