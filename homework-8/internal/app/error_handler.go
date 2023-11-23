package app

import (
	"context"
	"errors"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"google.golang.org/grpc"
)

type AppError interface {
	error
	StatusCode() int
}

func UnaryErrorHandlerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, err
		}

		var appErr AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		} else {
			return nil, &apperr.InternalServerError{Message: "Internal server error"}
		}
	}
}
