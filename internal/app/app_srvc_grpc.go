// Copyright(c) 2021 Altessa Solutions Inc.
// All rights reserved.

package app

import (
	"github.com/kitdoo/sn/internal/path"
	grpcService "github.com/kitdoo/sn/internal/transport/grpc"
	"github.com/kitdoo/sn/internal/transport/grpc/ping-pong"
	"github.com/kitdoo/sn/internal/transport/grpc/ping-pong/proto/pb"
	"google.golang.org/grpc"
)

func (a *App) startGRPCService() (err error) {
	defer func() {
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to start service")
			a.Shutdown()
		}
	}()

	var service *grpcService.Server
	if service, err = grpcService.New(a.config.GRPC, &grpcService.Interceptors{}, path.CertsCacheDir(), a.logger); err != nil {
		return
	}
	service.RegistrationService(func(gs *grpc.Server) {
		pingpongpb_v1.RegisterPingPongServer(gs, pingpong.New())
	})

	if err = service.Start(&a.servicesWait); err != nil {
		return
	}

	a.closers = append(a.closers, service)
	return
}
