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

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/mrcook/go-plugin-examples/bidirectional/proto"
)

// CounterStore is the interface that we're exposing as a plugin type.
// Any plugin that wishes to act as a CounterStore plugin must implement this interface.
type CounterStore interface {
	Put(key string, value int64, a AddHelper) error
	Get(key string) (int64, error)
}

// AddHelper is the interface being exposed for use with CounterStore plugins,
// which can call back to the host application for this Sum functionality.
type AddHelper interface {
	Sum(int64, int64) (int64, error)
}

// CounterPluginName is an important variable.
// All CounterStore plugins MUST use the same value in their plugin.ServeConfig
// when specifying the plugins (pluginMap).
// Its name can be any string value you wish.
const CounterPluginName = "counter"

// CounterPlugin is the implementation of plugin.Plugin used to serve and
// consume gRPC plugins of type CounterStore.
//
// Both host application and plugins use this implementation when they assign
// their gRPC plugins, i.e. via the `pluginMap`.
//
// NOTE: plugin.GRPCBroker is used to create a multiplexed stream on plugin
// connections so that both outgoing and incoming requests can be handled.
// The broker can remain unused when only host->plugin communication is needed.
type CounterPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	// Concrete implementation, written in Go.
	// This is only used for plugins that are written in Go.
	Impl CounterStore
}

// GRPCServer must return a gRPC server for this plugin type.
func (p *CounterPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCounterServer(s, &grpcCounterServer{Impl: p.Impl, broker: broker})
	return nil
}

// GRPCClient must return an implementation of our interface that communicates
// over a gRPC client.
func (p *CounterPlugin) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcCounterClient{client: proto.NewCounterClient(c), broker: broker}, nil
}
