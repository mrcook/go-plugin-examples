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
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
// Any plugin that wishes to act as a Greeter plugin must implement this interface.
type Greeter interface {
	Greet() string // TODO: add error on return? (see panic message in `Greet()` below)
}

// GreeterPluginName is an important variable.
// All Greeter plugins MUST use the same value in their plugin.ServeConfig when
// specifying the plugins (pluginMap).
// Its name can be any string value you wish.
const GreeterPluginName = "greeter"

// GreeterPlugin is the implementation of plugin.Plugin used to serve and
// consume net/rpc plugins of type Greeter.
//
// Both host applications and plugins use this implementation when they assign
// their plugins, i.e. via the `pluginMap`.
//
// NOTE: ignore MuxBroker. That is used to create more multiplexed streams on
// our plugin connection and is a more advanced use case.
type GreeterPlugin struct {
	Impl Greeter
}

// Server must return an RPC server for this plugin type. We construct a
// GreeterRpcServer for this.
func (p *GreeterPlugin) Server(_ *plugin.MuxBroker) (interface{}, error) {
	return &greeterServer{Impl: p.Impl}, nil
}

// Client must return an implementation of our interface that communicates over
// an RPC client. We return GreeterRpcClient for this.
func (_ *GreeterPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &greeterClient{client: c}, nil
}

// GreeterClient is a client implementation that talks over RPC.
type greeterClient struct{ client *rpc.Client }

func (g *greeterClient) Greet() string {
	var resp string

	// `Plugin`: a go-plugin hardcoded value
	// `Greet` the method as defined on the Greeter plugin interface
	err := g.client.Call("Plugin.Greet", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors,
		// if they don't, there isn't much other choice here.
		panic(err)
	}

	return resp
}

// GreeterServer is the RPC server that GreeterRpcClient talks to,
// conforming to the requirements of net/rpc.
type greeterServer struct {
	Impl Greeter
}

func (s *greeterServer) Greet(args interface{}, resp *string) error {
	*resp = s.Impl.Greet()
	return nil
}
