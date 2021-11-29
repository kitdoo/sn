package grpc

import (
	"net"
	"sync"
	"time"

	"github.com/kitdoo/sn/internal/atomic"

	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kitdoo/sn/internal/config"
	"github.com/kitdoo/sn/internal/transport/grpc/interceptors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	cfg          *config.GRPC
	grpc         *grpc.Server
	logger       zerolog.Logger
	started      atomic.Bool
	stopCh       chan struct{}
	certCacheDir string
	interceptors *Interceptors
	wait         *sync.WaitGroup
}

type Interceptors struct {
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

func New(cfg *config.GRPC, i *Interceptors, certCacheDir string, logger zerolog.Logger) (*Server, error) {
	s := &Server{
		cfg:          cfg,
		logger:       logger.With().Str("subsystem", "grpc").Logger(),
		stopCh:       make(chan struct{}),
		certCacheDir: certCacheDir,
		interceptors: i,
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		interceptors.UnaryServerReqIDInterceptor(),
		interceptors.UnaryServerLogInterceptor(s.logger),
	}

	if i != nil && len(i.UnaryInterceptors) > 0 {
		unaryInterceptors = append(unaryInterceptors, i.UnaryInterceptors...)
	}

	streamServerInterceptor := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		interceptors.StreamServerReqIDInterceptor(),
		interceptors.StreamServerLogInterceptor(s.logger),
	}

	if i != nil && len(i.StreamInterceptors) > 0 {
		streamServerInterceptor = append(streamServerInterceptor, i.StreamInterceptors...)
	}

	var grpcOptions = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamServerInterceptor...),
	}

	s.grpc = grpc.NewServer(grpcOptions...)
	return s, nil
}

func (s *Server) Start(wait *sync.WaitGroup) error {
	if !s.started.CompareAndSwap(false, true) {
		return nil
	}

	wait.Add(1)

	if s.cfg.Reflection {
		reflection.Register(s.grpc)
	}

	listen, err := net.Listen("tcp", s.cfg.ListenAddress)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to start service")
		return err
	}

	chErr := make(chan error, 1)
	go func() {
		defer func() {
			s.logger.Info().Msg("service was stopped")
			wait.Done()
		}()
		if err := s.grpc.Serve(listen); err != nil {
			s.logger.Error().Err(err).Msg("failed to start service")
			chErr <- err
		}
	}()

	//nolint:gomnd
	timer := time.NewTimer(time.Second * 2)
	defer timer.Stop()

	select {
	case err := <-chErr:
		return err
	case <-timer.C:
	}

	close(chErr)
	s.logger.Info().Msgf("service is listening on %v", listen.Addr())
	return nil
}

func (s *Server) RegistrationService(f func(gs *grpc.Server)) {
	f(s.grpc)
}

func (s *Server) Shutdown() error {
	if !s.started.CompareAndSwap(true, false) {
		return nil
	}

	if s.grpc != nil {
		s.grpc.GracefulStop()
	}

	return nil
}
