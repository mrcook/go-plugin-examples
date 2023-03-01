// An example using gRPC for communication between the host application and plugins.
//
// Three plugin examples are provided:
// - plugin-go-grpc and plugin-go-python, which communicate over gRPC
// - plugin-go-netrpc, which uses net/rpc
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/grpc/sdk"
)

// The executable for each plugin.
const (
	grpcPluginExecutable   = "./kv-go-grpc"
	rpcPluginExecutable    = "./kv-go-netrpc"
	pythonPluginExecutable = "python plugin-python/plugin.py"
)

func main() {
	// Fetch the plugin type, command, and key/value data from the CLI args.
	args := parseFlags()

	// Configure which plugin to use!
	var pluginName, pluginExecutable string
	if args.plugin == "grpc" {
		pluginName = sdk.KVStoreGrpcPluginName
		pluginExecutable = grpcPluginExecutable
	} else if args.plugin == "rpc" {
		pluginName = sdk.KVStoreNetRpcPluginName
		pluginExecutable = rpcPluginExecutable
	} else if args.plugin == "python" {
		pluginName = sdk.KVStoreGrpcPluginName
		pluginExecutable = pythonPluginExecutable
	}

	// PluginMap is the map of plugins we can dispense.
	pluginMap := map[string]plugin.Plugin{
		sdk.KVStoreGrpcPluginName:   &sdk.KVPluginGRPC{},
		sdk.KVStoreNetRpcPluginName: &sdk.KVPluginRPC{},
	}

	// Configure a new plugin client:
	// - HandshakeConfig: is required
	// - Plugins: is a map containing the supported plugins and their plugin.Plugin implementations
	// - Cmd: points to the compiled binary of your plugin
	// - AllowedProtocols: by default only net/rpc is allowed, so add gRPC support
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  sdk.HandshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command("sh", "-c", pluginExecutable),
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
	raw, err := client.Dispense(pluginName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// As Dispense() returns an interface, we need to cast it to the plugin
	// type supported by the host application, which in our case is a KV store.
	// This feels like a normal interface implementation, but is in fact
	// communicating over an RPC connection.
	kv := raw.(sdk.KVStore)

	// Call the appropriate method based on that requested by the user.
	if args.command == "get" {
		result, err := kv.Get(args.key)
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		// Let's see what the plugin returns!
		fmt.Println(string(result))
	} else if args.command == "put" {
		err := kv.Put(args.key, []byte(args.value))
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
	}
}

// Contains all the data required to run the application.
type cliArgs struct {
	plugin  string // the plugin to be used
	command string // get or put command
	key     string // custom key name (appended to the KV store filename)
	value   string // comment to be saved in the file
}

func parseFlags() cliArgs {
	grpc := flag.Bool("grpc", false, "App will use plugin-go-grpc.")
	rpc := flag.Bool("rpc", false, "App will use plugin-go-netrpc.")
	python := flag.Bool("python", false, "App will use plugin-python.")
	flag.Parse()

	var pluginType string
	if *grpc {
		pluginType = "grpc"
	} else if *rpc {
		pluginType = "rpc"
	} else if *python {
		pluginType = "python"
	} else {
		fmt.Println("a plugin type must be specified, include: --grpc, --rpc, or --python")
		os.Exit(1)
	}

	command := flag.Arg(0)
	if command != "get" && command != "put" {
		fmt.Printf("invalid command, must be 'get' or 'put', given '%s'\n", command)
		os.Exit(1)
	}

	key := flag.Arg(1)
	value := flag.Arg(2)
	if len(key) == 0 {
		fmt.Println("key must be present")
		os.Exit(1)
	} else if command == "put" && len(value) == 0 {
		fmt.Println("value must be provide with the 'put' command")
		os.Exit(1)
	}

	return cliArgs{
		plugin:  pluginType,
		command: command,
		key:     key,
		value:   value,
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
