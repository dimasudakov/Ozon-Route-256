package app

import (
	"context"
	logging "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		logging.SetGlobal(logger)

		ctx = logging.ToContext(ctx, logger)

		resp, err := handler(ctx, req)
		if err != nil {

		}
		return resp, err
	}
}
