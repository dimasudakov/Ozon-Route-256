package app

import (
	"errors"
	"net/http"
)

type AppError interface {
	error
	StatusCode() int
}

func ErrorHandler(f func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var appErr AppError
			if errors.As(err, &appErr) {
				http.Error(w, appErr.Error(), appErr.StatusCode())
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
	}
}
