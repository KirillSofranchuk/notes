package model

import (
	"fmt"
	"strings"
	"time"
)

type Folder struct {
	Id        int
	Title     string
	Timestamp time.Time
	UserId    int `json:"-"`
	Notes     []Note
}

func NewFolder(title string, userId int) (*Folder, *ApplicationError) {
	if len(title) == 0 {
		return nil, NewApplicationError(ErrorTypeValidation, "Название папки не может быть пустым", nil)
	}

	return &Folder{
		Id:        0,
		Title:     title,
		Timestamp: time.Now(),
		UserId:    userId,
	}, nil
}

func (f *Folder) GetInfo() string {
	notesInfo := "No notes"

	notes := f.GetNotes()

	if len(notes) > 0 {
		var notesList strings.Builder
		for _, note := range notes {
			notesList.WriteString(note.GetInfo())
		}
		notesInfo = fmt.Sprintf(notesList.String())
	}

	return fmt.Sprintf("Id: %d \n"+
		"Title: %s \n"+
		"TimeStamp: %s \n"+
		"UserId: %d \n"+
		"Notes: %s", f.Id, f.Title, f.Timestamp.Format(time.RFC1123), f.UserId, notesInfo)
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

func (f *Folder) AppendNotes(notes []Note) {
	if f.Notes == nil {
		f.Notes = notes
		return
	}

	f.Notes = append(f.Notes, notes...)
}
