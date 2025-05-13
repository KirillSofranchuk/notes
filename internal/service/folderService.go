package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
)

const (
	folderTitleIsNotFree = "Папка с таким же именем уже добавлена"
	fakeId               = -1
)

type AbstractFolderService interface {
	CreateFolder(userId int, title string) (int, *model.ApplicationError)
	UpdateFolder(userId int, folderId int, title string) *model.ApplicationError
	DeleteFolder(userId int, folderId int) *model.ApplicationError
}

type FolderService struct {
	repo repository.AbstractRepository
}

func NewConcreteFolderService(repository repository.AbstractRepository) AbstractFolderService {
	return &FolderService{
		repo: repository,
	}
}

func (f FolderService) CreateFolder(userId int, title string) (int, *model.ApplicationError) {
	folder, err := model.NewFolder(title, userId)

	if err != nil {
		return fakeId, err
	}

	if !f.isTitleIsFree(folder.Title) {
		return fakeId, model.NewApplicationError(model.ErrorTypeValidation, folderTitleIsNotFree, nil)
	}

	id, err := f.repo.SaveEntity(folder)

	if err != nil {
		return fakeId, err
	}

	return id, nil
}

func (f FolderService) UpdateFolder(userId int, folderId int, title string) *model.ApplicationError {
	folder, err := model.NewFolder(title, userId)

	if err != nil {
		return err
	}

	if !f.isTitleIsFree(folder.Title) {
		return model.NewApplicationError(model.ErrorTypeValidation, folderTitleIsNotFree, nil)
	}

	folderDb, err := f.repo.GetFolderById(folderId, userId)

	if err != nil {
		return err
	}

	folderDb.Title = title

	_, errSave := f.repo.SaveEntity(folderDb)
	return errSave
}

func (f FolderService) DeleteFolder(userId int, folderId int) *model.ApplicationError {
	folderDb, err := f.repo.GetFolderById(folderId, userId)

	if err != nil {
		if err.Type == model.ErrorTypeNotFound {
			return nil
		}
		return err
	}

	return f.repo.DeleteEntity(folderDb)
}

func (f FolderService) isTitleIsFree(title string) bool {
	folders := f.repo.GetFolders()

	for _, folder := range folders {
		if folder.Title == title {
			return false
		}
	}

	return true
}
