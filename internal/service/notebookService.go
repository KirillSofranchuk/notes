package service

import (
	"Notes/internal/model"
	"Notes/internal/repository"
)

type AbstractNotebookService interface {
	GetUserNotebook(userId int) model.Notebook
}

type ConcreteNotebookService struct {
	repo repository.AbstractRepository
}

func NewConcreteNotebookService(repository repository.AbstractRepository) AbstractNotebookService {
	return &ConcreteNotebookService{
		repo: repository,
	}
}

func (n *ConcreteNotebookService) GetUserNotebook(userId int) model.Notebook {
	folders := n.repo.GetFoldersByUserId(userId)
	notes := n.repo.GetNotesByUserId(userId)

	return model.Notebook{
		Folders: n.getFoldersWithNotes(folders, notes),
		Notes:   n.getNotesRelatedToFolder(notes, nil),
	}
}

func (n *ConcreteNotebookService) getFoldersWithNotes(folders []*model.Folder, notes []*model.Note) []model.Folder {
	userFolders := make([]model.Folder, 0)

	for _, folder := range folders {
		folderId := folder.Id
		relatedNotes := n.getNotesRelatedToFolder(notes, &folderId)
		folder.AppendNotes(relatedNotes)
		userFolders = append(userFolders, *folder)
	}

	return userFolders
}

func (n *ConcreteNotebookService) getNotesRelatedToFolder(notes []*model.Note, folderId *int) []model.Note {
	notesWithoutFolder := make([]model.Note, 0)

	for _, note := range notes {
		if note.GetFolderId() == folderId {
			notesWithoutFolder = append(notesWithoutFolder, *note)
		}
	}

	return notesWithoutFolder
}
