package model

import (
	"errors"
	"time"
)

type Folder struct {
	Id        int
	Title     string
	Timestamp time.Time
	UserId    int
	notes     *[]Note
}

func NewFolder(title string, userId int) (*Folder, error) {
	if len(title) == 0 {
		return nil, errors.New("folder name cannot be empty")
	}

	return &Folder{
		Id:        0,
		Title:     title,
		Timestamp: time.Now(),
		UserId:    userId,
	}, nil
}

func (f *Folder) GetNotes() *[]Note {
	if f.notes == nil {
		return &[]Note{}
	}

	return f.notes
}

func (f *Folder) AppendNote(note Note) {
	if f.notes == nil {
		f.notes = &[]Note{note}
		return
	}

	*f.notes = append(*f.notes, note)
}
