package logger

import (
	"Notes/internal/repository"
	"context"
	"log"
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

func LogEntitiesAsync(repo repository.AbstractRepository, ctx context.Context) {
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
						log.Printf("New user saved: %v \n", e.GetInfo())
						log.Println("_____________________________________")
					}

					previousUsersCount = len(users)
				}

				if len(folders) > previousFoldersCount {
					for _, e := range folders[previousFoldersCount:] {
						log.Printf("New folder saved: %v \n", e.GetInfo())
						log.Println("_____________________________________")
					}

					previousFoldersCount = len(folders)
				}

				if len(notes) > previousNotesCount {
					for _, e := range notes[previousNotesCount:] {
						log.Printf("New note saved: %v \n", e.GetInfo())
						log.Println("_____________________________________")
					}

					previousNotesCount = len(notes)
				}

			case <-ctx.Done():
				log.Println("logger stopped")
				return
			}
		}
	}()
}
