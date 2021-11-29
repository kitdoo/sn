package interceptors

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const RequestIDKey = "request-id"

//wrap srv stream

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func WrapServerStream(stream grpc.ServerStream, ctx context.Context) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	if ctx == nil {
		ctx = stream.Context()
	}
	return &wrappedServerStream{ServerStream: stream, ctx: ctx}
}

//ctx

type ctxRequestIDrMarker struct{}

var (
	ctxRequestIDKey = &ctxRequestIDrMarker{}
)

func newContextWithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ctxRequestIDKey, reqID)
}

func ReqIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(ctxRequestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

//addr

func ClientAddr(ctx context.Context) string {
	clientAddr := ""
	if p, ok := peer.FromContext(ctx); ok {
		clientAddr = fmt.Sprintf("%s/%s", p.Addr.Network(), p.Addr.String())
	}
	return clientAddr
}

// ctx for log
type grpcZerologLogger struct {
	logger zerolog.Logger
}

type ctxLoggerMarker struct{}

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

func newCtxWithLogger(ctx context.Context, log zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, &grpcZerologLogger{log})
}

func LoggerFromContext(ctx context.Context) zerolog.Logger {
	l, ok := ctx.Value(ctxLoggerKey).(*grpcZerologLogger)
	if !ok || l == nil {
		return zerolog.Nop()
	}
	return l.logger
}
