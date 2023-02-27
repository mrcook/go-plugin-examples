// A basic plugin example of type Greeter.
//
// Note: all logging requests from this plugin will be made via the host
// application itself. If the host application disables logging, then these
// logs will also be discarded.
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/basic/sdk"
)

// HelloGreeterPlugin is our custom plugin: it's a real implementation of the
// Greeter plugin type.
type HelloGreeterPlugin struct {
	logger hclog.Logger
}

// Greet is the message we wish to return from our custom plugin.
func (plugin *HelloGreeterPlugin) Greet() string {
	msg := "Hello!"

	// this log message will be sent to the host application.
	plugin.logger.Debug("HelloGreeterPlugin.Greet", "greeting", msg)

	return msg
}

// go-plugin's are normal Go applications so require a main entry point.
// Once the host application has loaded (dispensed) the plugin, go-plugin will
// start the plugin, and manage its full lifecycle.
func main() {
	greeter := &HelloGreeterPlugin{
		logger: logger(),
	}

	// this log message will be sent to the host application.
	greeter.logger.Debug("HelloGreeterPlugin main() function")

	// Assign our plugin as the required plugin type.
	var plugins = plugin.PluginSet{
		sdk.GreeterPluginName: &sdk.GreeterPlugin{Impl: greeter},
	}

	// start listening for incoming RPC requests.
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.HandshakeConfig,
		Plugins:         plugins,
	})
}

// Use the HashiCorp Logger for logging.
func logger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
}
