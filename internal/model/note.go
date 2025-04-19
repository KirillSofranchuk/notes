package model

import (
	"errors"
	"fmt"
	"time"
)

const MaxContentLength = 200
const MaxTagsCount = 3

type Note struct {
	Id         int
	Title      string
	Content    string
	UserId     int
	IsFavorite bool
	Timestamp  time.Time
	Tags       []string
	folder     *Folder
}

func NewNote(title string, content string, userId int, tags []string) (*Note, error) {
	validationError := validateNote(title, content, tags)

	if validationError != nil {
		return nil, validationError
	}

	return &Note{
		Id:         0,
		Title:      title,
		Content:    content,
		UserId:     userId,
		IsFavorite: false,
		Timestamp:  time.Now(),
		Tags:       tags,
		folder:     nil,
	}, nil
}

func (n Note) GetInfo() string {
	return fmt.Sprintf("Id: %d \n"+
		"Title: %s \n"+
		"Content: %s \n"+
		"UserId: %d \n"+
		"IsFavorite: %v \n"+
		"TimeStamp: %s \n"+
		"Tags: %v", n.Id, n.Title, n.Content, n.UserId, n.IsFavorite, n.Timestamp.Format(time.RFC1123), n.Tags)
}

func validateNote(title string, content string, tags []string) error {
	titleValidationError := validateTitle(title)
	if titleValidationError != nil {
		return titleValidationError
	}

	contentValidationError := validateContent(content)

	if contentValidationError != nil {
		return contentValidationError
	}

	tagsValidationError := validateTags(tags)

	if tagsValidationError != nil {
		return tagsValidationError
	}

	return nil
}

func validateTitle(title string) error {
	if len(title) == 0 {
		return errors.New("length cannot be empty")
	}

	return nil
}

func validateContent(content string) error {
	if len(content) == 0 {
		return errors.New("content cannot be empty")
	}

	if len(content) > MaxContentLength {
		return errors.New(fmt.Sprintf("content length cannot be more than %d symbols", MaxContentLength))
	}

	return nil
}

func validateTags(tags []string) error {
	if tags == nil {
		return nil
	}

	if len(tags) > MaxTagsCount {
		return errors.New(fmt.Sprintf("cannot add more than %d tags to note", MaxTagsCount))
	}

	return nil
}
