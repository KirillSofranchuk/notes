package model

import "time"

type FolderApi struct {
	Id        int
	Title     string
	Timestamp time.Time
	UserId    int `json:"-"`
	Notes     []NoteApi
}

func (f *FolderApi) AppendNotes(notes []NoteApi) {
	if f.Notes == nil {
		f.Notes = notes
		return
	}

	f.Notes = append(f.Notes, notes...)
}

func ToFolderApi(dbFolder *Folder) *FolderApi {
	if dbFolder == nil {
		return nil
	}

	return &FolderApi{
		Id:        dbFolder.Id,
		Title:     dbFolder.Title,
		Timestamp: dbFolder.Timestamp,
		UserId:    dbFolder.UserId,
	}
}

func ToFoldersApi(dbFolders []*Folder) []*FolderApi {
	folders := make([]*FolderApi, 0, len(dbFolders))
	for i := range dbFolders {
		folders = append(folders, &FolderApi{
			Id:        dbFolders[i].Id,
			Title:     dbFolders[i].Title,
			Timestamp: dbFolders[i].Timestamp,
			UserId:    dbFolders[i].UserId,
		})
	}

	return folders
}
