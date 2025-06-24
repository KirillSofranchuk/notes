package repository

import "Notes/internal/model"

//go:generate mockgen -source=abstractRepository.go -destination=../../internal/service/mock/abstractRepository.go -package=mock

type AbstractRepository interface {
	SaveEntity(entity model.BusinessEntity) (int, *model.ApplicationError)
	DeleteEntity(entity model.BusinessEntity) *model.ApplicationError
	GetUserById(id int) (*model.User, *model.ApplicationError)
	GetFolderById(id int, userId int) (*model.Folder, *model.ApplicationError)
	GetNoteById(id int, userId int) (*model.Note, *model.ApplicationError)
	GetUser(login, password string) (*model.User, *model.ApplicationError)
	GetFoldersByUserId(userId int) []*model.Folder
	GetNotesByUserId(userId int) []*model.Note
	GetUsers() []*model.User
}
