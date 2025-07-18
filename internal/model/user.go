package model

import (
	"fmt"
	"time"
	"unicode"
)

const MinLoginLength = 8
const MinPasswordLength = 10

type User struct {
	Id        int
	Name      string
	Surname   string
	Login     string
	Password  string
	Timestamp time.Time
}

func NewUser(name string, surname string, login string, password string) (*User, *ApplicationError) {
	validationError := validateUser(name, surname, login, password)
	if validationError != nil {
		return nil, validationError
	}

	return &User{
		Id:       0,
		Name:     name,
		Surname:  surname,
		Login:    login,
		Password: password,
	}, nil
}

func (u *User) SetId(id int) {
	u.Id = id
}

func (u *User) GetId() int {
	return u.Id
}

func (u *User) SetTimestamp() {
	u.Timestamp = time.Now()
}

func validateUser(name string, surname string, login string, password string) *ApplicationError {
	personalDataValidationError := validatePersonalData(name, surname)
	if personalDataValidationError != nil {
		return personalDataValidationError
	}

	credentialsValidationError := validateCredentials(login, password)

	if credentialsValidationError != nil {
		return credentialsValidationError
	}

	return nil
}

func validatePersonalData(name string, surname string) *ApplicationError {
	if len(name) == 0 {
		return NewApplicationError(ErrorTypeValidation, "Имя не может быть пустым", nil)
	}

	if len(surname) == 0 {
		return NewApplicationError(ErrorTypeValidation, "Фамилия не может быть пустой", nil)
	}

	return nil
}

func validateCredentials(login string, password string) *ApplicationError {
	loginValidation := validateLogin(login)
	if loginValidation != nil {
		return loginValidation
	}

	passwordValidation := validatePassword(password)
	if passwordValidation != nil {
		return passwordValidation
	}

	return nil
}

func validateLogin(login string) *ApplicationError {
	if len(login) < MinLoginLength {
		message := fmt.Sprintf("Логин слишком короткий. Пожалуйста, создайте логин длинной не меньше %d символов", MinLoginLength)
		return NewApplicationError(ErrorTypeValidation, message, nil)
	}

	return nil
}

func validatePassword(password string) *ApplicationError {
	if len(password) < MinPasswordLength {
		message := fmt.Sprintf("Пароль слишком короткий. Пожалуйста, создайте пароль длиной не менее %d символов", MinPasswordLength)
		return NewApplicationError(ErrorTypeValidation, message, nil)
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return NewApplicationError(ErrorTypeValidation, "Пароль должен содержать букву верхнего регистра.", nil)
	}
	if !hasLower {
		return NewApplicationError(ErrorTypeValidation, "Пароль должен содержать букву нижнего регистра.", nil)
	}
	if !hasNumber {
		return NewApplicationError(ErrorTypeValidation, "Пароль должен содержать число.", nil)
	}
	if !hasSpecial {
		return NewApplicationError(ErrorTypeValidation, "Пароль должен содержать спецсимволы.", nil)
	}

	return nil
}
