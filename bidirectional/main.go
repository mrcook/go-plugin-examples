// A gRPC example with two-way communication between the host application and plugin.
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
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/bidirectional/sdk"
)

func main() {
	// Fetch the command, key, and value from the CLI args.
	args := parseFlags()

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
		Cmd:              exec.Command("sh", "-c", "./counter-go-grpc"),
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
	if args.command == "get" {
		result, err := counter.Get(args.key)
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		// Let's see what the plugin returns!
		fmt.Println(result)
	} else if args.command == "put" {
		// Provide our trusted helper for doing the summation work.
		err = counter.Put(args.key, args.value, &hostAddHelper{})
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
	}
}

// As we're not trusting plugins to do the summation work, add our own
// helper that the RPC requests will use.
type hostAddHelper struct{}

func (*hostAddHelper) Sum(a, b int64) (int64, error) {
	return a + b, nil
}

// Contains all the data required to run the application.
type cliArgs struct {
	command string // get or put command
	key     string // filename key
	value   int64  // value to be added
}

func parseFlags() cliArgs {
	flag.Parse()

	command := flag.Arg(0)
	if command != "get" && command != "put" {
		fmt.Printf("invalid command, must be 'get' or 'put', given '%s'\n", command)
		os.Exit(1)
	}

	var numberToAdd int64

	key := flag.Arg(1)
	value := flag.Arg(2)
	if len(key) == 0 {
		fmt.Println("key must be present")
		os.Exit(1)
	} else if command == "put" {
		if i, err := strconv.Atoi(value); err != nil {
			fmt.Println("value does not seem to be a valid number:", err.Error())
			os.Exit(1)
		} else {
			numberToAdd = int64(i)
		}
	}

	return cliArgs{
		command: command,
		key:     key,
		value:   numberToAdd,
	}
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
