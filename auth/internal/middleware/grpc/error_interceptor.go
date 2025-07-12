package grpcmiddleware

import (
	"context"

	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ErrorHandlerInterceptor handles AppError and logs gRPC errors consistently.
func ErrorHandlerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// AppError handling
		if appErr, ok := err.(*errors.AppError); ok {
			logger.Error("Handled AppError",
				zap.String("method", info.FullMethod),
				zap.String("code", appErr.Code()),
				zap.String("public", appErr.Public()),
				zap.Error(appErr),
			)

			return nil, status.Errorf(
				errcode.MapCodeToGRPC(appErr.Code()),
				appErr.Public(),
			)
		}

		// Log unexpected errors
		st, _ := status.FromError(err)
		logger.Error("Unhandled gRPC error",
			zap.String("method", info.FullMethod),
			zap.String("grpc_code", st.Code().String()),
			zap.String("message", st.Message()),
			zap.Error(err),
		)

		return nil, status.Errorf(
			errcode.MapCodeToGRPC(errcode.ErrInternalFailure),
			"Internal server error",
		)
	}
}
