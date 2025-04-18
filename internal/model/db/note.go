package db

import "time"

type Note struct {
	Id         int
	Title      string
	Content    string
	UserId     int
	FolderId   int
	IsFavorite bool
	Timestamp  time.Time
	Tags       []string
}
