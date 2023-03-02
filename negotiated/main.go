// An example using a versioned plugin: one version uses net/rpc,
// while the other gRPC.
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

	"github.com/mrcook/go-plugin-examples/negotitated/sdk"
)

func main() {
	// Fetch the plugin version, command, key, and value from the CLI args.
	args := parseFlags()

	// Initialize the array of versioned plugins that can be dispensed.
	plugins := map[int]plugin.PluginSet{}

	// Both versions can be supported, but switch the implementation to
	// demonstrate version negotiation.
	if args.pluginVersion == 2 {
		plugins[2] = plugin.PluginSet{
			sdk.KVStorePluginName: &sdk.KVPluginRPC{},
		}
	} else if args.pluginVersion == 3 {
		plugins[3] = plugin.PluginSet{
			sdk.KVStorePluginName: &sdk.KVPluginGRPC{},
		}
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
	raw, err := client.Dispense(sdk.KVStorePluginName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// As Dispense() returns an interface, we need to cast it to the plugin
	// type supported by the host application, which in our case is a KVStore store.
	// This feels like a normal interface implementation, but is in fact
	// communicating over an RPC connection.
	kv := raw.(sdk.KVStore)

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
	pluginVersion int    // the plugin version to use
	command       string // get or put command
	key           string // custom key name (appended to the KV store filename)
	value         string // comment to be saved in the file
}

func parseFlags() cliArgs {
	pluginVersion := flag.Int("plugin", 3, "Plugin version to use: 2 (net/rpc) or 3 (gRPC)")
	flag.Parse()

	if *pluginVersion < 2 || *pluginVersion > 3 {
		fmt.Println("plugin version must be 2 or 3")
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
		pluginVersion: *pluginVersion,
		command:       command,
		key:           key,
		value:         value,
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
