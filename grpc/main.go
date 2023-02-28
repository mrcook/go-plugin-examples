// An example using gRPC for communication between the host application on plugins.
//
// Three plugin examples exist, two that communicate over gRPC (plugin-go-grpc,
// plugin-go-python) and one that uses net/rcp (plugin-go-netrpc)
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/grpc/sdk"
)

func main() {
	// PluginMap is the map of plugins we can dispense.
	pluginMap := map[string]plugin.Plugin{
		sdk.KVStoreGrpcPluginName:   &sdk.KVPluginGRPC{},
		sdk.KVStoreNetRpcPluginName: &sdk.KVPluginRPC{},
	}

	// Configure a new plugin client:
	// - HandshakeConfig: is required
	// - Plugins: is a map containing the supported plugins and their plugin.Plugin implementations
	// - Cmd: points to the compiled binary of your plugin (set via an ENV)
	// - AllowedProtocols: by default only net/rpc is allowed, so add gRPC support
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  sdk.HandshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("sh", "-c", os.Getenv("KV_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger:           logger(),
	})
	defer pluginClient.Kill()

	// Get the client for RPC communication.
	client, err := pluginClient.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Use the ENV to pick plugin to use.
	var pluginType string
	switch os.Getenv("KV_PLUGIN") {
	case "./kv-go-netrpc":
		pluginType = sdk.KVStoreNetRpcPluginName
	default:
		pluginType = sdk.KVStoreGrpcPluginName
	}

	// Request the plugin.
	raw, err := client.Dispense(pluginType)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// As Dispense() returns an interface, we need to cast it to the plugin
	// type supported by the host application, which in our case is a KV store.
	// This feels like a normal interface implementation, but is in fact
	// communicating over an net/rpc connection.
	kv := raw.(sdk.KVStore)
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "get":
		result, err := kv.Get(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		// Let's see what the plugin returns!
		fmt.Println(string(result))
	case "put":
		err := kv.Put(os.Args[1], []byte(os.Args[2]))
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("Please only use 'get' or 'put', given: %q", os.Args[0])
		os.Exit(1)
	}
}

// A HashiCorp Logger, configured to discard all logs.
func logger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: io.Discard,
		Level:  hclog.Debug,
	})
}
