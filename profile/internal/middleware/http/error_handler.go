package httpmiddleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()

			if err == nil {
				return
			}

			// Handler application errors
			if appErr, ok := err.Err.(*errors.AppError); ok {
				// Log the application error
				logger.Error("Handled AppError",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("path", ctx.FullPath()),
					zap.String("method", ctx.Request.Method),
					zap.String("code", appErr.Code()),
					zap.String("public", appErr.Public()),
					zap.Error(appErr),
				)

				// Capture request body
				ctx.JSON(errcode.MapCodeToHTTP(appErr.Code()), gin.H{
					"error": appErr.Public(),
					"code":  appErr.Code(),
				})
				return
			}

			// Log unexpected errors
			logger.Error("Unhandled internal error",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("path", ctx.FullPath()),
				zap.String("method", ctx.Request.Method),
				zap.Error(err.Err),
			)

			ctx.JSON(errcode.MapCodeToHTTP(errcode.ErrInternalFailure), gin.H{
				"error": "Internal server error",
				"code":  errcode.ErrInternalFailure,
			})
		}
	}
}
