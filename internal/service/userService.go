package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
	"Notes/internal/utils"
)

const loginAlreadyUsedMessage = "пользователь с таким логином уже добавлен"

type AbstractUserService interface {
	CreateUser(login, password, name, surname string) (int, *model.ApplicationError)
	UpdateUser(id int, login, password, name, surname string) *model.ApplicationError
	GetUser(userId int) (*model.User, *model.ApplicationError)
	DeleteUser(id int) *model.ApplicationError
}

type UserService struct {
	repo repository.AbstractRepository
}

func NewConcreteUserService(repository repository.AbstractRepository) AbstractUserService {
	return &UserService{
		repo: repository,
	}
}

func (u UserService) CreateUser(login, password, name, surname string) (int, *model.ApplicationError) {
	newUser, err := model.NewUser(name, surname, login, password)

	if err != nil {
		return -1, err
	}

	if !u.isLoginFree(newUser.Login, newUser.GetId()) {
		return -1, model.NewApplicationError(model.ErrorTypeValidation, loginAlreadyUsedMessage, nil)
	}

	passwordHash, errHash := utils.GetHash(newUser.Password)
	if errHash != nil {
		return -1, errHash
	}

	newUser.Password = passwordHash

	id, err := u.repo.SaveEntity(newUser)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (u UserService) UpdateUser(id int, login, password, name, surname string) *model.ApplicationError {
	newUser, err := model.NewUser(name, surname, login, password)

	if err != nil {
		return err
	}

	if !u.isLoginFree(newUser.Login, id) {
		return model.NewApplicationError(model.ErrorTypeValidation, loginAlreadyUsedMessage, nil)
	}

	userDb, err := u.repo.GetUserById(id)
	if err != nil {
		return err
	}

	passwordHash, errHash := utils.GetHash(newUser.Password)
	if errHash != nil {
		return errHash
	}

	userDb.Login = login
	userDb.Password = passwordHash
	userDb.Name = name
	userDb.Surname = surname

	_, err = u.repo.SaveEntity(userDb)
	if err != nil {
		return err
	}

	return nil
}

func (u UserService) GetUser(userId int) (*model.User, *model.ApplicationError) {
	user, err := u.repo.GetUserById(userId)
	return user, err
}

func (u UserService) DeleteUser(id int) *model.ApplicationError {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteEntity(user)
}

func (u UserService) isLoginFree(login string, userId int) bool {
	users := u.repo.GetUsers()

	for _, user := range users {
		if login == user.Login && userId != user.GetId() {
			return false
		}
	}

	return true
}
