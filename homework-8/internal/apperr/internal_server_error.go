package apperr

import "net/http"

type InternalServerError struct {
	Message string
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{
		Message: message,
	}
}

func (e InternalServerError) Error() string {
	return e.Message
}

func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}
