package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const filePermission = 0644

type BusinessEntityStorage[T BusinessEntity] struct {
	entities []T
	mu       sync.RWMutex
	filename string
}

func NewBusinessEntityStorage[T BusinessEntity](filename string) *BusinessEntityStorage[T] {
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

func (b *BusinessEntityStorage[T]) Save(entity T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entities = append(b.entities, entity)

	file, err := os.OpenFile(b.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, filePermission)
	if err != nil {
		return
	}
	defer file.Close()

	line, _ := json.Marshal(entity)

	_, _ = file.WriteString(string(line) + "\n")
	return
}

func (b *BusinessEntityStorage[T]) GetAll() []T {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.entities
}

func (b *BusinessEntityStorage[T]) GetCount() int {
	return len(b.entities)
}
