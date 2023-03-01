// A gRPC example with two-way communication between the host application and plugin.
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
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/bidirectional/sdk"
)

func main() {
	// A map of the plugins we can dispense.
	pluginMap := plugin.PluginSet{
		sdk.CounterPluginName: &sdk.CounterPlugin{},
	}

	// Configure a new plugin client:
	// - HandshakeConfig: is required
	// - Plugins: is a map containing the supported plugins and their plugin.Plugin implementations
	// - Cmd: points to the compiled binary of your plugin
	// - AllowedProtocols: by default only net/rpc is allowed, so add gRPC support
	// - Logger: configured to discard all logs
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  sdk.HandshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("sh", "-c", os.Getenv("COUNTER_PLUGIN")),
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
	raw, err := client.Dispense(sdk.CounterPluginName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// As Dispense() returns an interface, we need to cast it to the plugin
	// type supported by the host application, which in our case is a counter store.
	// This feels like a normal interface implementation, but is in fact
	// communicating over an RPC connection.
	counter := raw.(sdk.CounterStore)

	// Call the appropriate method based on that requested by the user.
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "get":
		result, err := counter.Get(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		// Let's see what the plugin returns!
		fmt.Println(result)
	case "put":
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		// Provide our trusted helper for doing the summation work.
		err = counter.Put(os.Args[1], int64(i), &hostAddHelper{})
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("Please only use 'get' or 'put'")
		os.Exit(1)
	}
}

// As we're not trusting plugins to do the summation work, have our own
// helper that the RPC requests will use.
type hostAddHelper struct{}

func (*hostAddHelper) Sum(a, b int64) (int64, error) {
	return a + b, nil
}

// A HashiCorp Logger, configured to discard all logs.
// If no logger is specified plugin.NewClient will use `hclog` by default.
func logger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: io.Discard, // use os.Stderr to show logs
		Level:  hclog.Debug,
	})
}
