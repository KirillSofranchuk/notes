package model

import (
	"errors"
	"fmt"
	"unicode"
)

const minLoginLength = 8
const minPasswordLength = 10

type User struct {
	Id       int
	Name     string
	Surname  string
	Login    string
	Password string
}

func NewUser(name string, surname string, login string, password string) (*User, error) {
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

func validateUser(name string, surname string, login string, password string) error {
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

func validatePersonalData(name string, surname string) error {
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}

	if len(surname) == 0 {
		return errors.New("surname cannot be empty")
	}

	return nil
}

func validateCredentials(login string, password string) error {
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

func validateLogin(login string) error {
	if len(login) < minLoginLength {
		return errors.New(fmt.Sprintf("Login is too short. Please create login with at least %d symbols length", minLoginLength))
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return errors.New(fmt.Sprintf("password is too short. Please create password with at least %d symbols length", minPasswordLength))
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
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
