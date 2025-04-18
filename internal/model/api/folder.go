package api

type FolderRequest struct {
	Title string
}

type FolderResponse struct {
	Id    int
	Title string
	Notes []NoteResponse
}
