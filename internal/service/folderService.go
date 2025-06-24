package service

//go:generate mockgen -source=folderService.go -destination=mock/folderService.go -package=mock

import (
	"Notes/internal/constants"
	"Notes/internal/model"
	"Notes/internal/repository"
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
		return constants.FakeId, err
	}

	if !f.isTitleIsFree(folder.Title, userId, 0) {
		return constants.FakeId, model.NewApplicationError(model.ErrorTypeValidation, constants.FolderTitleIsNotFree, nil)
	}

	id, err := f.repo.SaveEntity(folder)

	if err != nil {
		return constants.FakeId, err
	}

	return id, nil
}

func (f FolderService) UpdateFolder(userId int, folderId int, title string) *model.ApplicationError {
	folder, err := model.NewFolder(title, userId)

	if err != nil {
		return err
	}

	if !f.isTitleIsFree(folder.Title, userId, folderId) {
		return model.NewApplicationError(model.ErrorTypeValidation, constants.FolderTitleIsNotFree, nil)
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

func (f FolderService) isTitleIsFree(title string, userId int, folderId int) bool {
	folders := f.repo.GetFoldersByUserId(userId)
	for _, folder := range folders {
		if folder.Title == title && folder.Id != folderId {
			return false
		}
	}

	return true
}
