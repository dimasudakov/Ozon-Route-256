package app

import (
	"errors"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"net/http"
)

func ErrorHandler(f func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var appErr apperr.AppError
			if errors.As(err, &appErr) {
				http.Error(w, appErr.Error(), appErr.StatusCode())
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
	}
}
