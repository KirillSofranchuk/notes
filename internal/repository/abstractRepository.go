package repository

import "Notes/internal/model"

type AbstractRepository interface {
	SaveEntity(entity model.BusinessEntity) (int, *model.ApplicationError)
	DeleteEntity(entity model.BusinessEntity) *model.ApplicationError
	GetUsers() []*model.User
	GetNotes() []*model.Note
	GetFolders() []*model.Folder
	GetUserById(id int) (*model.User, *model.ApplicationError)
	GetFolderById(id int, userId int) (*model.Folder, *model.ApplicationError)
	GetNoteById(id int, userId int) (*model.Note, *model.ApplicationError)
	GetUsersCount() int
	GetNotesCount() int
	GetFoldersCount() int
	GetUser(login, password string) (*model.User, *model.ApplicationError)
	GetFoldersByUserId(userId int) []*model.Folder
	GetNotesByUserId(userId int) []*model.Note
}
