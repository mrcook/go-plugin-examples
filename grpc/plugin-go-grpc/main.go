// A plugin example of type KVStore, which communicates over gRPC.
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/grpc/sdk"
)

// the files for this plugin use the prefix:
const filenamePrefix = "kv_grpc_"

// GrpcPlugin is our custom plugin: it's a real implementation of the KVStore
// plugin type that writes to a local file with the key name and the contents
// are the value of the key.
type GrpcPlugin struct{}

// Put will overwrite the file contents with the new key/value data.
func (GrpcPlugin) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin-go-grpc", string(value)))
	return os.WriteFile(filenamePrefix+key, value, 0644)
}

// Get reads the file and returns the value stored for the matching key.
func (GrpcPlugin) Get(key string) ([]byte, error) {
	return os.ReadFile(filenamePrefix + key)
}

// go-plugin's are normal Go applications so require a main entry point.
// Once the host application has loaded (dispensed) the plugin, go-plugin will
// start the plugin, and manage its full lifecycle.
func main() {
	// Assign our plugin as the required plugin type.
	plugins := plugin.PluginSet{
		sdk.KVStoreGrpcPluginName: &sdk.KVPluginGRPC{Impl: &GrpcPlugin{}},
	}

	// start listening for incoming gRPC requests.
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.HandshakeConfig,
		Plugins:         plugins,

		// A non-nil value here enables gRPC serving for this plugin.
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
