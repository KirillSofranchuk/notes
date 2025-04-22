package repository

import (
	"Notes/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	usersFile      = "users.json"
	foldersFile    = "folders.json"
	notesFile      = "notes.json"
	filePermission = 0644
)

type PlainRepository struct {
	users     []*model.User
	muUsers   sync.RWMutex
	notes     []*model.Note
	muNotes   sync.RWMutex
	folders   []*model.Folder
	muFolders sync.RWMutex
}

func NewPlainRepository() AbstractRepository {
	return &PlainRepository{
		users:   make([]*model.User, 0),
		notes:   make([]*model.Note, 0),
		folders: make([]*model.Folder, 0),
	}
}

func (s *PlainRepository) LoadStoredData() {
	// Загружаем параллельно, так как файлы не зависят друг от друга
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		loadData(usersFile, &s.users)
	}()

	go func() {
		defer wg.Done()
		loadData(foldersFile, &s.folders)
	}()

	go func() {
		defer wg.Done()
		loadData(notesFile, &s.notes)
	}()

	wg.Wait()

	fmt.Println("Data loaded")
}

func (s *PlainRepository) SaveEntity(entity model.BusinessEntity) {
	switch e := entity.(type) {
	case *model.User:
		s.muUsers.Lock()
		s.users = append(s.users, e)
		appendToFile(usersFile, e)
		s.muUsers.Unlock()
	case *model.Note:
		s.muNotes.Lock()
		s.notes = append(s.notes, e)
		appendToFile(notesFile, e)
		s.muNotes.Unlock()
	case *model.Folder:
		s.muFolders.Lock()
		s.folders = append(s.folders, e)
		appendToFile(foldersFile, e)
		s.muFolders.Unlock()
	}
}

func (s *PlainRepository) GetUsers() []*model.User {
	s.muUsers.RLock()
	defer s.muUsers.RUnlock()
	return s.users
}

func (s *PlainRepository) GetNotes() []*model.Note {
	s.muNotes.RLock()
	defer s.muNotes.RUnlock()
	return s.notes
}

func (s *PlainRepository) GetFolders() []*model.Folder {
	s.muFolders.RLock()
	defer s.muFolders.RUnlock()
	return s.folders
}

func loadData[T model.BusinessEntity](filename string, receiver *[]T) {
	storedData := readEntities[T](filename)

	fmt.Printf("Loaded %d entities from %s \n", len(storedData), filename)

	*receiver = append(*receiver, storedData...)
}

func readEntities[T model.BusinessEntity](filename string) []T {
	entities := make([]T, 0)

	file, err := os.Open(filename)
	if err != nil {
		return entities
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entity T
		if err := json.Unmarshal(scanner.Bytes(), &entity); err != nil {
			fmt.Println()
		}
		entities = append(entities, entity)
	}

	return entities
}

func appendToFile[T model.BusinessEntity](filename string, item T) {
	// Не используем mutex так как чтение из файла только при старте программы,
	//а пишется за раз сущность только одного типа
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, filePermission)
	if err != nil {
		return
	}
	defer file.Close()

	line, _ := json.Marshal(item)

	_, _ = file.WriteString(string(line) + "\n")
	return
}
