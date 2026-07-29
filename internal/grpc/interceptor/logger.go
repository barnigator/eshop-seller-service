package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Logging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		response, err := handler(ctx, req)

		requestID := RequestIDFromContext(ctx)
		duration := time.Since(startedAt)
		code := status.Code(err)

		log.Log(
			ctx,
			logLevelByCode(code),
			"grpc request completed",
			slog.String("request_id", requestID),
			slog.String("method", info.FullMethod),
			slog.String("status", code.String()),
			slog.Duration("duration", duration),
		)

		return response, err
	}
}

func logLevelByCode(code codes.Code) slog.Level {
	switch code {
	case codes.OK:
		return slog.LevelInfo

	case codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.Unauthenticated,
		codes.PermissionDenied,
		codes.FailedPrecondition:
		return slog.LevelWarn

	case codes.Internal,
		codes.Unknown,
		codes.Unavailable,
		codes.DataLoss,
		codes.DeadlineExceeded:
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}
