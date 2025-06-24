package service

//go:generate mockgen -source=notebookService.go -destination=mock/notebookService.go -package=mock

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

	mappedNotes := model.ToNotesApi(notes)
	mappedFolders := model.ToFoldersApi(folders)

	return model.Notebook{
		Folders: n.getFoldersWithNotes(mappedFolders, mappedNotes),
		Notes:   n.getNotesRelatedToFolder(mappedNotes, nil),
	}
}

func (n *ConcreteNotebookService) getFoldersWithNotes(folders []*model.FolderApi, notes []*model.NoteApi) []model.FolderApi {
	userFolders := make([]model.FolderApi, 0)

	for _, folder := range folders {
		folderId := folder.Id
		relatedNotes := n.getNotesRelatedToFolder(notes, &folderId)
		folder.AppendNotes(relatedNotes)
		userFolders = append(userFolders, *folder)
	}

	return userFolders
}

func (n *ConcreteNotebookService) getNotesRelatedToFolder(notes []*model.NoteApi, folderId *int) []model.NoteApi {
	notesRelatedToFolder := make([]model.NoteApi, 0)

	for _, note := range notes {
		folderIdFromNote := note.FolderId
		if (folderIdFromNote == nil && folderId == nil) || ((folderIdFromNote != nil && folderId != nil) && (*folderIdFromNote == *folderId)) {
			notesRelatedToFolder = append(notesRelatedToFolder, *note)
		}
	}

	return notesRelatedToFolder
}
