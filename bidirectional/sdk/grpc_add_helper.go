// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/mrcook/go-plugin-examples/bidirectional/proto"
)

// grpcAddHelperClient is an implementation of AddHelper that talks over RPC.
type grpcAddHelperClient struct {
	client proto.AddHelperClient
}

func (c *grpcAddHelperClient) Sum(a, b int64) (int64, error) {
	resp, err := c.client.Sum(
		context.Background(),
		&proto.SumRequest{A: a, B: b},
	)
	if err != nil {
		hclog.Default().Info("add.Sum", "client", "start", "err", err)
		return 0, err
	}
	return resp.R, err
}

// grpcAddHelperServer is the gRPC server that grpcAddHelperClient talks to.
type grpcAddHelperServer struct {
	Impl AddHelper
}

func (s *grpcAddHelperServer) Sum(_ context.Context, req *proto.SumRequest) (*proto.SumResponse, error) {
	r, err := s.Impl.Sum(req.A, req.B)
	if err != nil {
		return nil, err
	}
	return &proto.SumResponse{R: r}, err
}
