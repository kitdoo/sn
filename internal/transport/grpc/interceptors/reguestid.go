package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerReqIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := uuid.New().String()
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ids := md[RequestIDKey]
			if len(ids) > 0 {
				requestID = ids[0]
			}
		}
		return handler(newContextWithReqID(ctx, requestID), req)
	}
}

func StreamServerReqIDInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		requestID := uuid.New().String()
		md, ok := metadata.FromIncomingContext(stream.Context())
		if ok {
			ids := md[RequestIDKey]
			if len(ids) > 0 {
				requestID = ids[0]
			}
		}
		return handler(srv, WrapServerStream(stream, newContextWithReqID(stream.Context(), requestID)))
	}
}
