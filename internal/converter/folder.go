package converter

import (
	"project/internal/model/api"
	"project/internal/model/db"
	"time"
)

func ApiFolderToDb(api api.FolderRequest, userId int) db.Folder {
	return db.Folder{
		Title:     api.Title,
		Timestamp: time.Now(),
		UserId:    userId,
		Notes:     nil,
	}
}

func DbFolderToApi(db db.Folder) api.FolderResponse {
	return api.FolderResponse{
		Id:    db.Id,
		Title: db.Title,
		Notes: getFolderNotes(db),
	}
}

func getFolderNotes(db db.Folder) []api.NoteResponse {
	if db.Notes == nil {
		return nil
	}

	var apiNotes []api.NoteResponse

	for _, note := range db.Notes {
		apiNotes = append(apiNotes, DbNoteToApi(note))
	}

	return apiNotes
}
