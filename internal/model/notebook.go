package model

// Notebook represents the notebook API response
// @Description Notebook information
type Notebook struct {
	Folders []FolderApi `json:"folders"`
	Notes   []NoteApi   `json:"notes"`
}
