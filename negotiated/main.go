// An example using a versioned plugin: one version uses net/rpc,
// while the other gRPC.
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

	"github.com/mrcook/go-plugin-examples/negotitated/sdk"
)

func main() {
	// Initialize an array of versioned plugins that can be dispensed.
	plugins := map[int]plugin.PluginSet{}

	// Both versions can be supported, but switch the implementation to
	// demonstrate version negotiation.
	switch os.Getenv("KV_PROTO") {
	case "netrpc":
		plugins[2] = plugin.PluginSet{
			sdk.KVStorePluginName: &sdk.KVPluginRPC{},
		}
	case "grpc":
		plugins[3] = plugin.PluginSet{
			sdk.KVStorePluginName: &sdk.KVPluginGRPC{},
		}
	default:
		fmt.Println("must set KV_PROTO to netrpc or grpc")
		os.Exit(1)
	}

	// Configure a new plugin client:
	// - HandshakeConfig: is required
	// - VersionedPlugins: is an array of plugin versions, and their plugin.Plugin implementations
	// - Cmd: points to the compiled binary of your plugin
	// - AllowedProtocols: by default only net/rpc is allowed, so add gRPC support
	// - Logger: configured to discard all logs
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  sdk.HandshakeConfig,
		VersionedPlugins: plugins,
		Cmd:              exec.Command("./kv-plugin"),
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

	// Request the plugin.
	raw, err := client.Dispense("kv")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// As Dispense() returns an interface, we need to cast it to the plugin
	// type supported by the host application, which in our case is a KVStore store.
	// This feels like a normal interface implementation, but is in fact
	// communicating over an RPC connection.
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
		fmt.Println("Please only use 'get' or 'put'")
		os.Exit(1)
	}
}

// A HashiCorp Logger, configured to discard all logs.
// If no logger is specified plugin.NewClient will use `hclog` by default.
func logger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: io.Discard,
		Level:  hclog.Debug,
	})
}
