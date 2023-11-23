package apperr

import "net/http"

type BadRequestError struct {
	Message string
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}
