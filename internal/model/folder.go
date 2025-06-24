package model

import (
	"time"
)

type Folder struct {
	Id        int
	Title     string
	Timestamp time.Time
	UserId    int
	Notes     []Note
}

func NewFolder(title string, userId int) (*Folder, *ApplicationError) {
	if len(title) == 0 {
		return nil, NewApplicationError(ErrorTypeValidation, "Название папки не может быть пустым", nil)
	}

	return &Folder{
		Id:     0,
		Title:  title,
		UserId: userId,
	}, nil
}

func (f *Folder) SetId(id int) {
	f.Id = id
}

func (f *Folder) GetId() int {
	return f.Id
}

func (f *Folder) GetNotes() []Note {
	if f.Notes == nil {
		return []Note{}
	}

	return f.Notes
}

func (f *Folder) SetTimestamp() {
	f.Timestamp = time.Now()
}
