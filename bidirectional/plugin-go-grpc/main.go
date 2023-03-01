// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/go-plugin"

	"github.com/mrcook/go-plugin-examples/bidirectional/sdk"
)

// The KV store filename prefix for this plugin.
const filenamePrefix = "kv_store_"

// CounterPlugin is our custom plugin: it's a real implementation of the
// CounterStore plugin type that updates and reads the number value stored
// in the local file.
type CounterPlugin struct{}

// storeData presents the JSON data stored in the local file.
type storeData struct {
	Value int64 `json:"value"`
}

// Put adds the given number to that stored in the file matching the key.
// Before writing the data an RPC request is made to the host application
// using the sdk.AddHelper.
func (k *CounterPlugin) Put(key string, value int64, adder sdk.AddHelper) error {
	v, _ := k.Get(key)

	// Request the host application to add the two numbers. This feels like a
	// normal method call but is in fact over an RPC connection.
	r, err := adder.Sum(v, value)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(&storeData{r})
	if err != nil {
		return err
	}

	return os.WriteFile(filenamePrefix+key, buf, 0644)
}

// Get reads the file for matching the key and returns the value stored therein.
func (k *CounterPlugin) Get(key string) (int64, error) {
	fileContents, err := os.ReadFile(filenamePrefix + key)
	if err != nil {
		return 0, err
	}

	data := &storeData{}
	err = json.Unmarshal(fileContents, data)
	if err != nil {
		return 0, err
	}

	return data.Value, nil
}

// go-plugin's are normal Go applications so require a main entry point.
// Once the host application has loaded (dispensed) the plugin, go-plugin will
// start the plugin, and manage its full lifecycle.
func main() {
	// Assign our plugin as the required plugin type.
	plugins := plugin.PluginSet{
		sdk.CounterPluginName: &sdk.CounterPlugin{Impl: &CounterPlugin{}},
	}

	// start listening for incoming gRPC requests.
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.HandshakeConfig,
		Plugins:         plugins,

		// A non-nil value here enables gRPC serving for this plugin.
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
