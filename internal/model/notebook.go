package model

// Notebook represents the notebook API response
// @Description Notebook information
type Notebook struct {
	Folders []Folder `json:"folders"`
	Notes   []Note   `json:"notes"`
}
