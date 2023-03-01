// Package sdk contains shared data between the host and plugins.
//
// This can be thought of as a Plugin Software Development Kit that plugin
// authors will use when developing their plugins for use within a specific
// host application.
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sdk

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mrcook/go-plugin-examples/negotitated/proto"
)

// KVStore is the interface that we're exposing as a plugin.
// Any plugin that wishes to act as a KVStore plugin must implement this interface.
type KVStore interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

// KVStorePluginName is an important variable.
// All CounterStore plugins MUST use the same value in their plugin.ServeConfig
// when specifying the plugins (pluginMap).
// Its name can be any string value you wish.
const KVStorePluginName = "kv"

// KVPluginRPC is the implementation of plugin.Plugin used to serve and consume
// net/rpc plugins of type KVStore.
//
// Both host applications and plugins use this implementation when they assign
// their net/rpc plugins, i.e. via the `pluginMap`.
//
// NOTE: ignore MuxBroker. That is used to create more multiplexed streams on
// our plugin connection and is a more advanced use case.
type KVPluginRPC struct {
	// Concrete implementation, written in Go.
	// This is only used for plugins that are written in Go.
	Impl KVStore
}

// Server must return an RPC server for this plugin type.
func (p *KVPluginRPC) Server(_ *plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

// Client must return an implementation of our interface that communicates over
// an RPC client.
func (_ *KVPluginRPC) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &rpcClient{client: c}, nil
}

// KVPluginGRPC is the implementation of plugin.GRPCPlugin used to serve and
// consume gRPC plugins of type KVStore.
//
// Both host applications and plugins use this implementation when they assign
// their gRPC plugins, i.e. via the `pluginMap`.
//
// NOTE: ignore GRPCBroker. That is used to create more multiplexed streams on
// our plugin connection and is a more advanced use case, e.g. bi-directional
// communication.
type KVPluginGRPC struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin

	// Concrete implementation, written in Go.
	// This is only used for plugins that are written in Go.
	Impl KVStore
}

// GRPCServer must return a gRPC server for this plugin type.
func (p *KVPluginGRPC) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterKVServer(s, &grpcServer{Impl: p.Impl})
	return nil
}

// GRPCClient must return an implementation of our interface that communicates
// over a gRPC client.
func (p *KVPluginGRPC) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: proto.NewKVClient(c)}, nil
}
