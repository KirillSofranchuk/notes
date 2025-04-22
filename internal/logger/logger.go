package logger

import (
	"Notes/internal/repository"
	"fmt"
	"time"
)

const milliSecondsToSleep = 200

func LogEntitiesAsync(repo repository.AbstractRepository, stop <-chan struct{}) {
	go func() {
		previousUsersCount := 0
		previousNotesCount := 0
		previousFoldersCount := 0

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
