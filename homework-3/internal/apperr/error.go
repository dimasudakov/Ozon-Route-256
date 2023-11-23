package apperr

type AppError interface {
	error
	StatusCode() int
}
