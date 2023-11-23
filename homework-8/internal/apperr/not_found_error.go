package apperr

import "net/http"

type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}
