// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"golang.org/x/net/context"

	"github.com/mrcook/go-plugin-examples/negotitated/proto"
)

// GRPCClient is an implementation of KVStore that talks over RPC.
type grpcClient struct{ client proto.KVClient }

func (m *grpcClient) Put(key string, value []byte) error {
	_, err := m.client.Put(context.Background(), &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *grpcClient) Get(key string) ([]byte, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// GRPCServer is the gRPC server that GRPCClient talks to.
type grpcServer struct {
	Impl KVStore
}

func (m *grpcServer) Put(_ context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *grpcServer) Get(_ context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}
