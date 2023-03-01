# Negotiated HashiCorp go-plugin Example

This example builds a simple key/value store CLI where the plugin version can
be negotiated between client and server.

The provided plugin supports two different versions, one which communicates
over gRPC, and one that uses net/rpc.

## Negotiated Protocol

The Client sends the list of available plugin versions to the server. When
presented with a list of plugin versions, the server iterates over them in
reverse and uses the highest numbered match to choose the plugin to execute.
If a legacy client is used and no versions are sent to the server, the server
will default to the oldest version in its configuration.


## Usage

A `Makefile` is provide for ease of use. Running `make` will compile both the
host application and plugin, and then run the application with some example
requests.

Additional make commands:

```sh
make build  # build the app and plugin binaries
make pbuf   # re-generate the protocol buffers
make clean  # remove all binaries and store files.
```

The application accepts two commands: `get` and `put`. The `put` command takes
two arguments: a _key_ and a string _value_. The key will be appended to the
filename, while the value will be saved to that file.

Here's a full example:

```sh
# Write a value using proto version 3 and gRPC
KV_PROTO=grpc ./app put hello "Planet Earth"

# Read it back using proto version 2 and net/rpc
KV_PROTO=netrpc ./app get hello
Planet Earth

Written from plugin version 3
Read by plugin version 2
```


## LICENSE

All new code and documentation, copyright (c) 2023 Michael R. Cook.

Based on the examples from https://github.com/hashicorp/go-plugin, copyright (c) HashiCorp, Inc.

SPDX-License-Identifier: MPL-2.0
