package repository

import (
	"Notes/internal/model"
	"fmt"
	"sync"
)

const (
	usersFile      = "users.json"
	foldersFile    = "folders.json"
	notesFile      = "notes.json"
	filePermission = 0644
)

type PlainRepository struct {
	users   *model.BusinessEntityStorage[*model.User]
	notes   *model.BusinessEntityStorage[*model.Note]
	folders *model.BusinessEntityStorage[*model.Folder]
}

func NewPlainRepository() AbstractRepository {
	repo := &PlainRepository{
		users:   model.NewBusinessEntityStorage[*model.User](usersFile),
		notes:   model.NewBusinessEntityStorage[*model.Note](notesFile),
		folders: model.NewBusinessEntityStorage[*model.Folder](foldersFile),
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		repo.users.Load()
	}()

	go func() {
		defer wg.Done()
		repo.folders.Load()
	}()

	go func() {
		defer wg.Done()
		repo.notes.Load()
	}()

	wg.Wait()

	fmt.Println("Data loaded")

	return repo
}

func (s *PlainRepository) SaveEntity(entity model.BusinessEntity) {
	switch e := entity.(type) {
	case *model.User:
		s.users.Save(e)
	case *model.Note:
		s.notes.Save(e)
	case *model.Folder:
		s.folders.Save(e)
	}
}

func (s *PlainRepository) GetUsers() []*model.User {
	return s.users.GetAll()
}

func (s *PlainRepository) GetNotes() []*model.Note {
	return s.notes.GetAll()
}

func (s *PlainRepository) GetFolders() []*model.Folder {
	return s.folders.GetAll()
}

func (s *PlainRepository) GetUsersCount() int {
	return s.users.GetCount()
}

func (s *PlainRepository) GetFoldersCount() int {
	return s.folders.GetCount()
}

func (s *PlainRepository) GetNotesCount() int {
	return s.notes.GetCount()
}
