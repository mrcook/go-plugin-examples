// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"golang.org/x/net/context"

	"github.com/mrcook/go-plugin-examples/grpc/proto"
)

// grpcClient is an implementation of KVStore that talks over RPC.
type grpcClient struct {
	client proto.KVClient
}

func (c *grpcClient) Put(key string, value []byte) error {
	_, err := c.client.Put(context.Background(), &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (c *grpcClient) Get(key string) ([]byte, error) {
	resp, err := c.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// grpcServer is the gRPC server that grpcClient talks to.
type grpcServer struct {
	Impl KVStore
}

func (s *grpcServer) Put(_ context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.Put(req.Key, req.Value)
}

func (s *grpcServer) Get(_ context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := s.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}
