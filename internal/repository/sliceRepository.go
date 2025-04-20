package repository

import (
	"Notes/internal/model"
	"fmt"
)

var (
	users   []*model.User
	notes   []*model.Note
	folders []*model.Folder
)

type SliceRepository struct {
}

func NewSliceRepository() AbstractRepository {
	return &SliceRepository{}
}

func (s *SliceRepository) SaveEntity(entity model.BusinessEntity) {
	switch e := entity.(type) {
	case *model.User:
		users = append(users, e)
		fmt.Printf("New user saved: %v \n", e.GetInfo())
		fmt.Println("_____________________________________")
	case *model.Note:
		notes = append(notes, e)
		fmt.Printf("New note saved: %v \n", e.GetInfo())
		fmt.Println("_____________________________________")
	case *model.Folder:
		folders = append(folders, e)
		fmt.Printf("New folder saved: %v \n", e.GetInfo())
		fmt.Println("_____________________________________")
	}
}
