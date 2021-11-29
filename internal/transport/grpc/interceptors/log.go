package interceptors

import (
	"path"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"context"
)

const (
	KindField = "span.kind"
)

const (
	fieldNameClientAddr   = "client.addr"
	fieldNameClientApp    = "client.app"
	fieldNameGRPCCode     = "grpc.code"
	fieldNameGRPCDuration = "grpc.duration"
	fieldNameRequestID    = "request.id"
)

func UnaryServerLogInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, logger, info.FullMethod, startTime)

		resp, err := handler(newCtx, req)

		statusCode := status.Code(err)
		logLevel := statusCodeToLogLevel(statusCode)

		ll := LoggerFromContext(newCtx).With().
			Str(fieldNameClientAddr, ClientAddr(newCtx)).
			Str(fieldNameGRPCCode, statusCode.String()).
			Str(fieldNameGRPCDuration, time.Since(startTime).String())

		if id := ReqIDFromContext(newCtx); id != "" {
			ll = ll.Str(fieldNameRequestID, id)
		}

		if md, ok := metadata.FromIncomingContext(newCtx); ok {
			if userAgent, ok := md["user-agent"]; ok {
				ll = ll.Str(fieldNameClientApp, userAgent[0])
			}
		}

		if err != nil {
			ll = ll.Err(err)
		}

		doLog(ll.Logger(), logLevel, "finished unary call")
		return resp, err
	}
}

func StreamServerLogInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		newCtx := newLoggerForCall(stream.Context(), logger, info.FullMethod, startTime)

		streamContext := WrapServerStream(stream, newCtx)

		err := handler(srv, streamContext)

		statusCode := status.Code(err)
		logLevel := statusCodeToLogLevel(statusCode)

		ll := LoggerFromContext(newCtx).With().
			Str(fieldNameClientAddr, ClientAddr(stream.Context())).
			Str(fieldNameGRPCCode, statusCode.String()).
			Str(fieldNameGRPCDuration, time.Since(startTime).String())

		if id := ReqIDFromContext(newCtx); id != "" {
			ll = ll.Str(fieldNameRequestID, id)
		}

		if md, ok := metadata.FromIncomingContext(newCtx); ok {
			if userAgent, ok := md["user-agent"]; ok {
				ll = ll.Str(fieldNameClientApp, userAgent[0])
			}
		}

		if err != nil {
			ll = ll.Err(err)
		}

		doLog(ll.Logger(), logLevel, "finished stream call")
		return err
	}
}

//nolint:exhaustive
func doLog(logger zerolog.Logger, level zerolog.Level, msg string) {
	switch level {
	case zerolog.DebugLevel:
		logger.Debug().Msg(msg)
	case zerolog.InfoLevel:
		logger.Info().Msg(msg)
	case zerolog.WarnLevel:
		logger.Warn().Msg(msg)
	case zerolog.ErrorLevel:
		logger.Error().Msg(msg)
	case zerolog.FatalLevel:
		logger.Fatal().Msg(msg)
	case zerolog.PanicLevel:
		logger.Panic().Msg(msg)
	}
}

func newLoggerForCall(ctx context.Context, logger zerolog.Logger, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)

	ll := logger.With().
		Str(KindField, "server").
		Str("grpc.service", service).
		Str("grpc.method", method).
		Str("grpc.start_time", start.Format(time.RFC3339))

	if d, ok := ctx.Deadline(); ok {
		ll = ll.Str("grpc.request.deadline", d.Format(time.RFC3339))
	}

	return newCtxWithLogger(ctx, ll.Logger())
}

func statusCodeToLogLevel(code codes.Code) zerolog.Level {
	switch code {
	case codes.OK, codes.Canceled:
		return zerolog.InfoLevel
	case codes.Unknown, codes.InvalidArgument, codes.Unauthenticated,
		codes.FailedPrecondition, codes.PermissionDenied, codes.Unimplemented,
		codes.Internal, codes.Unavailable, codes.DataLoss:
		return zerolog.ErrorLevel
	case codes.DeadlineExceeded, codes.NotFound,
		codes.AlreadyExists, codes.ResourceExhausted, codes.Aborted,
		codes.OutOfRange:
		return zerolog.WarnLevel
	default:
		return zerolog.ErrorLevel
	}
}
