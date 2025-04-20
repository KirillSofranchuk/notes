package main

import (
	"Notes/internal/model"
	"fmt"
)

func main() {
	user, errUser := model.NewUser("Name", "Surname", "login1234", "Password1234$")
	if errUser != nil {
		fmt.Println("Error create new user. Reason: ", errUser)
		return
	}

	fmt.Println("New user created")
	fmt.Println(user.GetInfo())
	fmt.Println("_______________________________________________")

	folder, errFolder := model.NewFolder("title", 1)

	if errFolder != nil {
		fmt.Println("Error create folder. Reason: ", errUser)
		return
	}

	fmt.Println("New folder created")
	fmt.Println(folder.GetInfo())
	fmt.Println("_______________________________________________")

	noteWithoutTags, errNoteWithoutTags := model.NewNote("Title", "Sample text", 1, nil)

	if errNoteWithoutTags != nil {
		fmt.Println("Error create note without tags. Reason: ", errNoteWithoutTags)
		return
	}

	fmt.Println("New note without tags created")
	fmt.Println(noteWithoutTags.GetInfo())
	fmt.Println("_______________________________________________")

	noteWithTags, errNoteWithTags := model.NewNote("Title2", "Sample text2", 1, []string{"tag1", "tag2"})

	if errNoteWithTags != nil {
		fmt.Println("Error create note without tags. Reason: ", errNoteWithoutTags)
		return
	}

	fmt.Println("New note with tags created")
	fmt.Println(noteWithTags.GetInfo())
	fmt.Println("_______________________________________________")

	folder.AppendNote(*noteWithoutTags)

	fmt.Println("Note without tags added to folder")
	fmt.Println(folder.GetInfo())
	fmt.Println("_______________________________________________")
}
