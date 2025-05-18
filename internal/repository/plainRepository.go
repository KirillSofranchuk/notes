package repository

import (
	"Notes/internal/model"
	"Notes/internal/utils"
	"fmt"
	"sync"
)

const (
	usersFile   = "users.json"
	foldersFile = "folders.json"
	notesFile   = "notes.json"
)

type PlainRepository struct {
	users   *BusinessEntityStorage[*model.User]
	notes   *BusinessEntityStorage[*model.Note]
	folders *BusinessEntityStorage[*model.Folder]
}

func NewPlainRepository() AbstractRepository {
	repo := &PlainRepository{
		users:   NewBusinessEntityStorage[*model.User](usersFile),
		notes:   NewBusinessEntityStorage[*model.Note](notesFile),
		folders: NewBusinessEntityStorage[*model.Folder](foldersFile),
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

func (s *PlainRepository) SaveEntity(entity model.BusinessEntity) (int, *model.ApplicationError) {
	switch e := entity.(type) {
	case *model.User:
		return s.users.Save(e)
	case *model.Note:
		return s.notes.Save(e)
	case *model.Folder:
		return s.folders.Save(e)
	}

	return -1, nil
}

func (s *PlainRepository) DeleteEntity(entity model.BusinessEntity) *model.ApplicationError {
	switch e := entity.(type) {
	case *model.User:
		return s.users.Delete(e, true)
	case *model.Note:
		return s.notes.Delete(e, true)
	case *model.Folder:
		return s.folders.Delete(e, true)
	}

	return nil
}

func (s *PlainRepository) GetUserById(id int) (*model.User, *model.ApplicationError) {
	user, err := s.users.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PlainRepository) GetFolderById(id int, userId int) (*model.Folder, *model.ApplicationError) {
	folder, err := s.folders.GetById(id)
	if err != nil {
		return nil, err
	}

	if folder.UserId != userId {
		return nil, model.NewApplicationError(model.ErrorTypeNotFound, "Папка не найдена", nil)
	}

	return folder, nil
}

func (s *PlainRepository) GetNoteById(id int, userId int) (*model.Note, *model.ApplicationError) {
	note, err := s.notes.GetById(id)
	if err != nil {
		return nil, err
	}

	if note.UserId != userId {
		return nil, model.NewApplicationError(model.ErrorTypeNotFound, "Заметка не найдена", nil)
	}
	return note, nil
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

func (s *PlainRepository) GetUser(login, password string) (*model.User, *model.ApplicationError) {
	users := s.users.GetAll()

	for _, user := range users {

		if user.Login == login {
			arePasswordsEqual, err := utils.CompareHashAndPassword(user.Password, password)

			if err != nil {
				return nil, err
			}

			if arePasswordsEqual {
				return user, nil
			}
		}
	}

	return nil, model.NewApplicationError(model.ErrorTypeNotFound, "Пользователь не найден", nil)
}

func (s *PlainRepository) GetFoldersByUserId(userId int) []*model.Folder {
	folders := s.folders.GetAll()

	userFolders := make([]*model.Folder, 0)

	for _, folder := range folders {
		if folder.UserId == userId {
			userFolders = append(userFolders, folder)
		}
	}

	return userFolders
}

func (s *PlainRepository) GetNotesByUserId(userId int) []*model.Note {
	notes := s.notes.GetAll()

	userNotes := make([]*model.Note, 0)

	for _, note := range notes {
		if note.UserId == userId {
			userNotes = append(userNotes, note)
		}
	}

	return userNotes
}
