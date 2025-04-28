package logger

import (
	"Notes/internal/repository"
	"fmt"
	"time"
)

const milliSecondsToSleep = 200

var previousUsersCount = 0
var previousNotesCount = 0
var previousFoldersCount = 0

func InitLogger(repo repository.AbstractRepository) {
	previousNotesCount = repo.GetNotesCount()
	previousFoldersCount = repo.GetFoldersCount()
	previousUsersCount = repo.GetUsersCount()
}

func LogEntitiesAsync(repo repository.AbstractRepository, stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(milliSecondsToSleep * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				users := repo.GetUsers()
				folders := repo.GetFolders()
				notes := repo.GetNotes()

				if len(users) > previousUsersCount {
					for _, e := range users[previousUsersCount:] {
						fmt.Printf("New user saved: %v \n", e.GetInfo())
						fmt.Println("_____________________________________")
					}

					previousUsersCount = len(users)
				}

				if len(folders) > previousFoldersCount {
					for _, e := range folders[previousFoldersCount:] {
						fmt.Printf("New folder saved: %v \n", e.GetInfo())
						fmt.Println("_____________________________________")
					}

					previousFoldersCount = len(folders)
				}

				if len(notes) > previousNotesCount {
					for _, e := range notes[previousNotesCount:] {
						fmt.Printf("New note saved: %v \n", e.GetInfo())
						fmt.Println("_____________________________________")
					}

					previousNotesCount = len(notes)
				}

			case <-stop:
				return
			}
		}
	}()
}
