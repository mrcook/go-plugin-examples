// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"net/rpc"
)

// rpcClient is an implementation of KVStore that talks over RPC.
type rpcClient struct {
	client *rpc.Client
}

func (m *rpcClient) Put(key string, value []byte) error {
	// We don't expect a response, so we can just use interface{}
	var resp interface{}

	// The args are just going to be a map. A struct could be better.
	//
	// `Plugin`: a go-plugin hardcoded value
	// `Put` the method as defined on the KVStore plugin interface
	return m.client.Call(
		"Plugin.Put",
		map[string]interface{}{"key": key, "value": value},
		&resp,
	)
}

func (m *rpcClient) Get(key string) ([]byte, error) {
	var resp []byte

	// `Plugin`: a go-plugin hardcoded value
	// `Get` the method as defined on the KVStore plugin interface
	err := m.client.Call("Plugin.Get", key, &resp)

	return resp, err
}

// rpcServer is the RPC server that rpcClient talks to, conforming to
// the requirements of net/rpc
type rpcServer struct {
	Impl KVStore
}

func (m *rpcServer) Put(args map[string]interface{}, resp *interface{}) error {
	return m.Impl.Put(args["key"].(string), args["value"].([]byte))
}

func (m *rpcServer) Get(key string, resp *[]byte) error {
	v, err := m.Impl.Get(key)
	*resp = v
	return err
}
