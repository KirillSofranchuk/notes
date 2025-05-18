package model

import "fmt"

type ErrorType string

const (
	ErrorTypeDatabase   ErrorType = "DATABASE_ERROR"
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
	ErrorTypeAuth       ErrorType = "AUTH_ERROR"
)

type ApplicationError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (a *ApplicationError) Error() string {
	if a.Err != nil {
		return fmt.Sprintf("%s: %v", a.Message, a.Err)
	}

	return a.Message
}

func (a *ApplicationError) Unwrap() error {
	return a.Err
}

func NewApplicationError(errorType ErrorType, message string, err error) *ApplicationError {
	return &ApplicationError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

type ApiError struct {
	Message string
	Err     error
	Code    int
}

func (a *ApiError) Error() string {
	if a.Err != nil {
		return fmt.Sprintf("%s: %v", a.Message, a.Err)
	}

	return a.Message
}

func (a *ApiError) Unwrap() error {
	return a.Err
}

func newApiError(code int, message string, err error) *ApiError {
	return &ApiError{
		Message: message,
		Err:     err,
		Code:    code,
	}
}

func GetAppropriateApiError(appError *ApplicationError) *ApiError {
	switch appError.Type {
	case ErrorTypeDatabase:
	case ErrorTypeInternal:
		return newApiError(500, appError.Message, appError.Err)
	case ErrorTypeValidation:
		return newApiError(400, appError.Message, appError.Err)
	case ErrorTypeNotFound:
		return newApiError(404, appError.Message, appError.Err)
	case ErrorTypeAuth:
		return newApiError(400, appError.Message, appError.Err)
	}

	return newApiError(500, "Ошибка сервера", nil)
}
