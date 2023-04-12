// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mrcook/go-plugin-examples/bidirectional/proto"
)

// grpcCounterClient is an implementation of CounterStore that talks over RPC.
type grpcCounterClient struct {
	broker *plugin.GRPCBroker
	client proto.CounterClient
}

func (c *grpcCounterClient) Put(key string, value int64, a AddHelper) error {
	addHelperServer := &grpcAddHelperServer{Impl: a}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterAddHelperServer(s, addHelperServer)

		return s
	}

	brokerID := c.broker.NextId()
	go c.broker.AcceptAndServe(brokerID, serverFunc)

	_, err := c.client.Put(context.Background(), &proto.PutRequest{
		AddServer: brokerID,
		Key:       key,
		Value:     value,
	})

	s.Stop()
	return err
}

func (c *grpcCounterClient) Get(key string) (int64, error) {
	resp, err := c.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return 0, err
	}

	return resp.Value, nil
}

// grpcCounterServer is the gRPC server that grpcCounterClient talks to.
type grpcCounterServer struct {
	proto.UnimplementedCounterServer // enable forward-compatibility

	Impl CounterStore

	broker *plugin.GRPCBroker
}

func (s *grpcCounterServer) Put(_ context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	conn, err := s.broker.Dial(req.AddServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := &grpcAddHelperClient{proto.NewAddHelperClient(conn)}
	return &proto.Empty{}, s.Impl.Put(req.Key, req.Value, a)
}

func (s *grpcCounterServer) Get(_ context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := s.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}
