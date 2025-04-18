package converter

import (
	"project/internal/model/api"
	"project/internal/model/db"
	"time"
)

func ApiNoteToDb(api api.NoteRequest, userId int) db.Note {
	return db.Note{
		Content:   api.Content,
		Tags:      api.Tags,
		Title:     api.Title,
		UserId:    userId,
		Timestamp: time.Now(),
	}
}

func DbNoteToApi(db db.Note) api.NoteResponse {
	return api.NoteResponse{
		Id:         db.Id,
		Title:      db.Title,
		Content:    db.Content,
		Timestamp:  db.Timestamp,
		IsFavorite: db.IsFavorite,
		Tags:       db.Tags,
	}
}
