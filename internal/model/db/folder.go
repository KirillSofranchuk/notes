package db

import "time"

type Folder struct {
	Id        int
	Title     string
	Timestamp time.Time
	UserId    int
	Notes     []Note
}
