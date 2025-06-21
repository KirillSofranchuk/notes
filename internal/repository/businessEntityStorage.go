package repository

import (
	"Notes/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	filePermission = 0644
	defaultId      = 0
	defaultIndex   = -1
)

var (
	EntityNotFoundError = model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil)
	DataBaseError       = model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil)
)

type BusinessEntityStorage[T model.BusinessEntity] struct {
	entities []T
	mu       sync.RWMutex
	filename string
}

func NewBusinessEntityStorage[T model.BusinessEntity](filename string) *BusinessEntityStorage[T] {
	return &BusinessEntityStorage[T]{
		entities: make([]T, 0),
		filename: filename,
	}
}

func (b *BusinessEntityStorage[T]) Load() {
	file, err := os.Open(b.filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entity T
		if err := json.Unmarshal(scanner.Bytes(), &entity); err != nil {
			fmt.Println()
		}
		b.entities = append(b.entities, entity)
	}

	fmt.Printf("Loaded %d entities from %s\n", len(b.entities), b.filename)
}

func (b *BusinessEntityStorage[T]) Save(entity T) (int, *model.ApplicationError) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entity.SetTimestamp()

	if entity.GetId() == defaultId {
		newId := b.getIdForNewEntity()
		entity.SetId(newId)
		return b.saveNewEntity(entity)
	}

	return b.updateEntity(entity)
}

func (b *BusinessEntityStorage[T]) GetAll() []T {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.entities
}

func (b *BusinessEntityStorage[T]) GetById(id int) (T, *model.ApplicationError) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for i := range b.entities {
		if b.entities[i].GetId() == id {
			return b.entities[i], nil
		}
	}

	return *new(T), EntityNotFoundError
}

func (b *BusinessEntityStorage[T]) GetCount() int {
	return len(b.entities)
}

func (b *BusinessEntityStorage[T]) Delete(entity T, needToLock bool) *model.ApplicationError {
	if needToLock {
		b.mu.Lock()
		defer b.mu.Unlock()
	}

	sliceIndex := b.getIndexInSlice(entity)
	b.deleteFromSlice(sliceIndex)

	fileIndex, err := b.getIndexInFile(entity)

	if err != nil {
		return err
	}

	return b.deleteFromFile(fileIndex)
}

func (b *BusinessEntityStorage[T]) getIdForNewEntity() int {
	maxId := 1

	for _, item := range b.entities {
		if item.GetId() > maxId {
			maxId = item.GetId()
		}
	}

	return maxId
}

func (b *BusinessEntityStorage[T]) deleteFromSlice(index int) {
	if len(b.entities) == 1 {
		b.entities = make([]T, 0)
		return
	}
	b.entities = append(b.entities[:index], b.entities[index+1:]...)
}

func (b *BusinessEntityStorage[T]) deleteFromFile(index int) *model.ApplicationError {
	file, err := os.Open(b.filename)
	if err != nil {
		return DataBaseError
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		if currentLine != index {
			lines = append(lines, scanner.Text()+"\n")
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return DataBaseError
	}

	output := strings.Join(lines, "\n")

	writeErr := os.WriteFile(b.filename, []byte(output), 0644)

	if writeErr != nil {
		return DataBaseError
	}
	return nil
}

func (b *BusinessEntityStorage[T]) getIndexInSlice(entity T) int {
	for _, item := range b.entities {
		if item.GetId() == entity.GetId() {
			return item.GetId()
		}
	}

	return defaultIndex
}

func (b *BusinessEntityStorage[T]) getIndexInFile(entity T) (int, *model.ApplicationError) {
	file, err := os.Open(b.filename)
	if err != nil {
		return defaultIndex, DataBaseError
	}
	defer file.Close()

	index := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entityFromFile T
		if err := json.Unmarshal(scanner.Bytes(), &entityFromFile); err != nil {
			return defaultIndex, DataBaseError
		}

		if entityFromFile.GetId() == entity.GetId() {
			return index, nil
		}

		index++
	}

	return -1, DataBaseError
}

func (b *BusinessEntityStorage[T]) saveNewEntity(entity T) (int, *model.ApplicationError) {
	b.entities = append(b.entities, entity)

	file, err := os.OpenFile(b.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, filePermission)
	if err != nil {
		return -1, DataBaseError
	}
	defer file.Close()

	line, _ := json.Marshal(entity)

	_, _ = file.WriteString(string(line) + "\n")
	return entity.GetId(), nil
}

func (b *BusinessEntityStorage[T]) updateEntity(entity T) (int, *model.ApplicationError) {
	err := b.Delete(entity, false)
	if err != nil {
		return entity.GetId(), err
	}

	return b.saveNewEntity(entity)
}
