package api

import "time"

type NoteRequest struct {
	Title   string
	Content string
	Tags    []string
}

type NoteResponse struct {
	Id         int
	Title      string
	Content    string
	Timestamp  time.Time
	IsFavorite bool
	Tags       []string
}
