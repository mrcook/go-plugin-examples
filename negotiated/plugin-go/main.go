// A plugin example of type KVStore, which contains two different versions:
// one that communicates over gRPC, the other over net/rpc.
//
// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/negotitated/sdk"
)

// The KV store filename prefix for this plugin.
const filenamePrefix = "kv_store_"

// GrpcPlugin is v3 of our custom plugin: it's a real implementation of the
// KVStore plugin type that writes to a local file with the key name and the
// contents are the value of the key.
// It communicates with the host application via gRPC.
type GrpcPlugin struct{}

// Put will overwrite the file contents with the new key/value data.
// When the file is written the plugin version number will be appended.
func (GrpcPlugin) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin version 3\n", string(value)))
	return os.WriteFile(filenamePrefix+key, value, 0644)
}

// Get reads the file and returns the value stored for the matching key.
// Before returning the file contents, the plugin version number is appended.
func (GrpcPlugin) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(filenamePrefix + key)
	if err != nil {
		return nil, err
	}
	return append(d, []byte("Read by plugin version 3\n")...), nil
}

// NetRpcPlugin is v2 of our custom plugin: it's a real implementation of the
// KVStore plugin type that writes to a local file with the key name and the
// contents are the value of the key.
// It communicates with the host application via net/rpc.
type NetRpcPlugin struct{}

// Put will overwrite the file contents with the new key/value data.
// When the file is written the plugin version number will be appended.
func (NetRpcPlugin) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin version 2\n", string(value)))
	return os.WriteFile(filenamePrefix+key, value, 0644)
}

// Get reads the file and returns the value stored for the matching key.
// Before returning the file contents, the plugin version number is appended.
func (NetRpcPlugin) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(filenamePrefix + key)
	if err != nil {
		return nil, err
	}
	return append(d, []byte("Read by plugin version 2\n")...), nil
}

// go-plugin's are normal Go applications so require a main entry point.
// Once the host application has loaded (dispensed) the plugin, go-plugin will
// start the plugin, and manage its full lifecycle.
func main() {
	// Assign the version to the required plugin type.
	// - version 2 uses NetRPC
	// - version 3 uses GRPC
	versionedPlugins := map[int]plugin.PluginSet{
		2: {sdk.KVStorePluginName: &sdk.KVPluginRPC{Impl: &NetRpcPlugin{}}},
		3: {sdk.KVStorePluginName: &sdk.KVPluginGRPC{Impl: &GrpcPlugin{}}},
	}

	// start listening for incoming gRPC requests.
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig:  sdk.HandshakeConfig,
		VersionedPlugins: versionedPlugins,

		// A non-nil value here enables gRPC serving for this plugin.
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
