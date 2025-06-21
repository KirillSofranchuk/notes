package utils

import (
	"Notes/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hash1, password string) (bool, *model.ApplicationError) {
	err := bcrypt.CompareHashAndPassword([]byte(hash1), []byte(password))

	if err != nil {
		return false, model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при проверке хэша", err)
	}

	return true, nil
}
