// Copyright(c) 2021 Altessa Solutions Inc.
// All rights reserved.

package pingpong

import (
	"context"

	"github.com/kitdoo/sn/internal/transport/grpc/ping-pong/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pingpongpb_v1.UnimplementedPingPongServer
}

func New() *Service {
	return &Service{}
}

func (s *Service) Ping(ctx context.Context, in *emptypb.Empty) (*pingpongpb_v1.PingResponse, error) {
	return &pingpongpb_v1.PingResponse{Pong: "pong"}, nil
}
