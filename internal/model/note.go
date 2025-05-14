package model

import (
	"fmt"
	"time"
)

const MaxContentLength = 200
const MaxTagsCount = 3

type Note struct {
	Id         int
	Title      string
	Content    string
	UserId     int `json:"-"`
	IsFavorite bool
	Timestamp  time.Time
	Tags       *[]string
	folderId   *int
}

func NewNote(title string, content string, userId int, tags *[]string) (*Note, *ApplicationError) {
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
		folderId:   nil,
	}, nil
}

func (n *Note) GetInfo() string {
	return fmt.Sprintf("Id: %d \n"+
		"Title: %s \n"+
		"Content: %s \n"+
		"UserId: %d \n"+
		"IsFavorite: %v \n"+
		"TimeStamp: %s \n"+
		"Tags: %v", n.Id, n.Title, n.Content, n.UserId, n.IsFavorite, n.Timestamp.Format(time.RFC1123), n.Tags)
}

func (n *Note) SetId(id int) {
	n.Id = id
}

func (n *Note) GetId() int {
	return n.Id
}

func (n *Note) GetFolderId() *int {
	return n.folderId
}

func (n *Note) SetFolderId(folderId *int) {
	n.folderId = folderId
}

func validateNote(title string, content string, tags *[]string) *ApplicationError {
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

func validateTitle(title string) *ApplicationError {
	if len(title) == 0 {
		return NewApplicationError(ErrorTypeValidation, "Название заметки не может быть пустым", nil)
	}

	return nil
}

func validateContent(content string) *ApplicationError {
	if len(content) == 0 {
		return NewApplicationError(ErrorTypeValidation, "Заметка не может быть пустой", nil)
	}

	if len(content) > MaxContentLength {
		message := fmt.Sprintf("Длина заметки не может превышать %d символов", MaxContentLength)
		return NewApplicationError(ErrorTypeValidation, message, nil)
	}

	return nil
}

func validateTags(tags *[]string) *ApplicationError {
	if tags == nil {
		return nil
	}

	if len(*tags) > MaxTagsCount {
		message := fmt.Sprintf("Нельзя добавить больше, чем %d тегов к заметке.", MaxTagsCount)
		return NewApplicationError(ErrorTypeValidation, message, nil)
	}

	return nil
}
