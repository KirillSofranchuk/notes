package repository

import (
	"Notes/internal/model"
	"sync"
)

type SliceRepository struct {
	users     []*model.User
	muUsers   sync.RWMutex
	notes     []*model.Note
	muNotes   sync.RWMutex
	folders   []*model.Folder
	muFolders sync.RWMutex
}

func NewSliceRepository() AbstractRepository {
	return &SliceRepository{
		users:   make([]*model.User, 0),
		notes:   make([]*model.Note, 0),
		folders: make([]*model.Folder, 0),
	}
}

func (s *SliceRepository) SaveEntity(entity model.BusinessEntity) {
	switch e := entity.(type) {
	case *model.User:
		s.muUsers.Lock()
		s.users = append(s.users, e)
		s.muUsers.Unlock()
	case *model.Note:
		s.muNotes.Lock()
		s.notes = append(s.notes, e)
		s.muNotes.Unlock()
	case *model.Folder:
		s.muFolders.Lock()
		s.folders = append(s.folders, e)
		s.muFolders.Unlock()
	}
}

func (s *SliceRepository) GetUsers() []*model.User {
	s.muUsers.RLock()
	defer s.muUsers.RUnlock()
	return s.users
}

func (s *SliceRepository) GetNotes() []*model.Note {
	s.muNotes.RLock()
	defer s.muNotes.RUnlock()
	return s.notes
}

func (s *SliceRepository) GetFolders() []*model.Folder {
	s.muFolders.RLock()
	defer s.muFolders.RUnlock()
	return s.folders
}
