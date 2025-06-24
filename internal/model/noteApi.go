package model

import (
	"time"
)

type NoteApi struct {
	Id         int
	Title      string
	Content    string
	UserId     int `json:"-"`
	IsFavorite bool
	Timestamp  time.Time
	Tags       []string
	FolderId   *int `json:"-"`
}

func ToNoteApi(dbNote *Note) *NoteApi {
	if dbNote == nil {
		return nil
	}
	return &NoteApi{
		Id:         dbNote.Id,
		Title:      dbNote.Title,
		Content:    dbNote.Content,
		IsFavorite: dbNote.IsFavorite,
		Timestamp:  dbNote.Timestamp,
		Tags:       dbNote.Tags,
		FolderId:   dbNote.FolderId,
	}
}

func ToNotesApi(dbNotes []*Note) []*NoteApi {
	notes := make([]*NoteApi, 0, len(dbNotes))
	for i := range dbNotes {
		notes = append(notes, &NoteApi{
			Id:         dbNotes[i].Id,
			Title:      dbNotes[i].Title,
			Content:    dbNotes[i].Content,
			IsFavorite: dbNotes[i].IsFavorite,
			Timestamp:  dbNotes[i].Timestamp,
			Tags:       dbNotes[i].Tags,
			FolderId:   dbNotes[i].FolderId,
		})
	}
	return notes
}
