package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const requestIDHeader = "x-request-id"

type contextKey string

const requestIDContextKey contextKey = "request_id"

func RequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestID := requestIDFromMetadata(ctx)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx = context.WithValue(ctx, requestIDContextKey, requestID)

		return handler(ctx, req)
	}
}

func requestIDFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(requestIDHeader)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey).(string)
	return requestID
}
