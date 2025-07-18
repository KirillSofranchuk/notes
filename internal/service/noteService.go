package service

import (
	"Notes/internal/constants"
	"Notes/internal/model"
	"Notes/internal/repository"
	"strings"
)

//go:generate mockgen -source=noteService.go -destination=mock/noteService.go -package=mock

type AbstractNoteService interface {
	CreateNote(userId int, title string, content string, tags *[]string) (int, *model.ApplicationError)
	DeleteNote(userId int, id int) *model.ApplicationError
	UpdateNote(userId int, id int, title string, content string, tags *[]string) *model.ApplicationError
	MoveToFolder(userId int, id int, folderId *int) *model.ApplicationError
	AddToFavorites(userId int, id int) *model.ApplicationError
	DeleteFromFavorites(userId int, id int) *model.ApplicationError
	FindNotesByQueryPhrase(userId int, query string) []*model.NoteApi
	GetFavoriteNotes(userId int) []*model.NoteApi
}

type NoteService struct {
	repo repository.AbstractRepository
}

func NewConcreteNoteService(repository repository.AbstractRepository) AbstractNoteService {
	return &NoteService{
		repo: repository,
	}
}

func (n *NoteService) CreateNote(userId int, title string, content string, tags *[]string) (int, *model.ApplicationError) {
	newNote, err := model.NewNote(title, content, userId, tags)

	if err != nil {
		return constants.FakeId, err
	}

	if !n.isTitleFree(newNote.Title, userId, 0) {
		return constants.FakeId, model.NewApplicationError(model.ErrorTypeValidation, constants.NoteNameIsNotFree, nil)
	}

	return n.repo.SaveEntity(newNote)
}

func (n *NoteService) DeleteNote(userId int, id int) *model.ApplicationError {
	note, err := n.repo.GetNoteById(id, userId)

	if err != nil {
		if err.Type == model.ErrorTypeNotFound {
			return nil
		}

		return err
	}

	return n.repo.DeleteEntity(note)
}

func (n *NoteService) UpdateNote(userId int, id int, title string, content string, tags *[]string) *model.ApplicationError {
	noteModel, err := model.NewNote(title, content, userId, tags)

	if err != nil {
		return err
	}

	if !n.isTitleFree(title, userId, id) {
		return model.NewApplicationError(model.ErrorTypeValidation, constants.NoteNameIsNotFree, nil)
	}

	noteDb, err := n.repo.GetNoteById(id, userId)

	if err != nil {
		return err
	}

	noteDb.Title = noteModel.Title
	noteDb.Content = noteModel.Content
	noteDb.Tags = noteModel.Tags

	_, saveErr := n.repo.SaveEntity(noteDb)
	return saveErr
}

func (n *NoteService) MoveToFolder(userId int, id int, folderId *int) *model.ApplicationError {
	note, err := n.repo.GetNoteById(id, userId)

	if err != nil {
		return err
	}

	if folderId != nil {
		_, folderErr := n.repo.GetFolderById(*folderId, userId)

		if folderErr != nil {
			return folderErr
		}
	}

	note.FolderId = folderId

	_, errSave := n.repo.SaveEntity(note)

	if errSave != nil {
		return errSave
	}

	return nil
}

func (n *NoteService) AddToFavorites(userId int, id int) *model.ApplicationError {
	note, err := n.repo.GetNoteById(id, userId)

	if err != nil {
		return err
	}

	note.IsFavorite = true

	_, errSave := n.repo.SaveEntity(note)

	if errSave != nil {
		return errSave
	}

	return nil
}

func (n *NoteService) DeleteFromFavorites(userId int, id int) *model.ApplicationError {
	note, err := n.repo.GetNoteById(id, userId)

	if err != nil {
		return err
	}

	note.IsFavorite = false

	_, errSave := n.repo.SaveEntity(note)

	if errSave != nil {
		return errSave
	}

	return nil
}

func (n *NoteService) FindNotesByQueryPhrase(userId int, query string) []*model.NoteApi {
	userNotes := n.repo.GetNotesByUserId(userId)

	if query == "" {
		return model.ToNotesApi(userNotes)
	}

	relatedNotes := make([]*model.NoteApi, 0)

	for _, note := range userNotes {
		if strings.Contains(note.Title, query) || strings.Contains(note.Content, query) || n.containsTag(note.Tags, query) {
			relatedNotes = append(relatedNotes, model.ToNoteApi(note))
		}
	}

	return relatedNotes
}

func (n *NoteService) GetFavoriteNotes(userId int) []*model.NoteApi {
	userNotes := n.repo.GetNotesByUserId(userId)

	favoriteNotes := make([]*model.NoteApi, 0)

	for _, note := range userNotes {
		if note.IsFavorite {
			favoriteNotes = append(favoriteNotes, model.ToNoteApi(note))
		}
	}

	return favoriteNotes
}

func (n *NoteService) containsTag(tags []string, tag string) bool {
	if tags == nil {
		return false
	}

	for _, item := range tags {
		if item == tag {
			return true
		}
	}

	return false
}

func (n *NoteService) isTitleFree(title string, userId int, noteId int) bool {
	notes := n.repo.GetNotesByUserId(userId)

	for _, note := range notes {
		if note.Title == title && note.Id != noteId {
			return false
		}
	}

	return true
}
