package util

import (
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultGRPCErrorMessage = "Internal Server Error"
)

// InvalidArgumentError is used for gRPC invalid argument errors
var InvalidArgumentError = status.Error(codes.InvalidArgument, "Invalid Argument")

// NotFoundError is used for gRPC not found errors
var NotFoundError = status.Error(codes.NotFound, "Not Found")

// PermissionDeniedError is used for gRPC permission denied errors
var PermissionDeniedError = status.Error(codes.PermissionDenied, "Permission Denied")

func NewGRPCError(err error) error {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*errors.AppError); ok {
		return status.Errorf(
			errcode.MapCodeToGRPC(appErr.Code()),
			appErr.Public(),
		)
	}

	return status.Errorf(
		codes.Internal,
		DefaultGRPCErrorMessage,
	)
}
