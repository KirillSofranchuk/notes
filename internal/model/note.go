package model

import (
	"Notes/internal/constants"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Note struct {
	Id         int
	Title      string
	Content    string
	UserId     int
	IsFavorite bool
	Timestamp  time.Time
	Tags       pq.StringArray `gorm:"type:text[]"`
	FolderId   *int
}

func (n *Note) SetId(id int) {
	n.Id = id
}

func (n *Note) GetId() int {
	return n.Id
}

func (n *Note) SetTimestamp() {
	n.Timestamp = time.Now()
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
		Tags:       getTags(tags),
	}, nil
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

	if len(content) > constants.MaxContentLength {
		message := fmt.Sprintf("Длина заметки не может превышать %d символов", constants.MaxContentLength)
		return NewApplicationError(ErrorTypeValidation, message, nil)
	}

	return nil
}

func validateTags(tags *[]string) *ApplicationError {
	if tags == nil {
		return nil
	}

	if len(*tags) > constants.MaxTagsCount {
		message := fmt.Sprintf("Нельзя добавить больше, чем %d тегов к заметке.", constants.MaxTagsCount)
		return NewApplicationError(ErrorTypeValidation, message, nil)
	}

	return nil
}

func getTags(tags *[]string) []string {
	if tags == nil {
		return make([]string, 0)
	}

	return *tags
}
