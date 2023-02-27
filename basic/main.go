// A simple host application which will communicate over net/rpc with the basic
// "greeter" plugin example, as provided in the "hello_plugin_example"
// directory.
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/basic/sdk"
)

func main() {
	// The set of plugins that our host application supports.
	pluginMap := plugin.PluginSet{
		sdk.GreeterPluginName: &sdk.GreeterPlugin{},
	}

	// Configure a new plugin client:
	// - HandshakeConfig: is required
	// - Plugins: is a map containing the name of your plugin and its plugin.Plugin implementation
	// - Cmd: points to the compiled binary of your plugin
	// - Logger: (optional) used for logging from both the host application and your plugin (if configured)
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: sdk.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./hello_plugin"),
		Logger:          logger(),
	})
	defer pluginClient.Kill()

	// Get the client for RPC communication.
	client, err := pluginClient.Client()
	if err != nil {
		log.Fatal(err)
	}

	// In essence, load the plugin.
	raw, err := client.Dispense(sdk.GreeterPluginName)
	if err != nil {
		log.Fatal(err)
	}

	// As Dispense() returns an interface, we need to cast it to the plugin
	// type supported by the host application, which in our case is a Greeter.
	// This feels like a normal interface implementation, but is in fact
	// communicating over an net/rpc connection.
	greeter := raw.(sdk.Greeter)

	// Let's see what greeting the plugin returns!
	greeting := greeter.Greet()
	fmt.Printf("\n\nThe plugin greeting is: %s\n\n\n", greeting)
}

// Use the HashiCorp Logger for all our logs.
func logger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr, // io.Discard to omit logs
		Level:  hclog.Debug,
	})
}
