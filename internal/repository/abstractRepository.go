package repository

import "Notes/internal/model"

type AbstractRepository interface {
	SaveEntity(entity model.BusinessEntity)
	GetUsers() []*model.User
	GetNotes() []*model.Note
	GetFolders() []*model.Folder
	LoadStoredData()
}
