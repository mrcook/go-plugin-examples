# gRPC HashiCorp go-plugin Example

This example builds a simple key/value store CLI where the mechanism
for storing and retrieving keys is pluggable.

Three plugin implementations are provided: two that communicate over gRPC, and
one that uses net/rpc.

## main.go

Is the "host" application code.

A new client is created that can communicate over both gRPC and net/rpc.
Depending on the CLI flag set, one of the three plugins is loaded and a request
made to its `Get` or `Put` method.

When `Get` is called, the contents of the `kv_` file is printed to the terminal.

## plugin-*

The `plugin-go-grpc` and `plugin-python` plugins both communicate over gRPC,
while `plugin-go-netrpc` communicates over net/rpc.

You will need Python installed on your system to run the `plugin-python` example.

## sdk

The Plugin Software Development Kit.

This package (which could be released as a standalone package) is used by
plugin authors for creating `KVStore` plugins. It provides the `KVPluginGRPC`,
which is an implementation of the `plugin.GRPCPlugin`, and `KVPluginRPC`, which
is an implementation of the `plugin.Plugin`. These are used by `go-plugin` to
provide communication between the host application and plugins using either
gRPC or net/rpc.

### proto

Contains the gPRC protocol buffer definitions used by the SDK.


## Usage

A `Makefile` is provide for ease of use. Running `make` will compile the host
application and the two Go plugins.

Additional make commands:

```sh
make build-go  # build all Go binaries
make pbufs     # re-generate all protocol buffers
make clean     # remove all binaries and kv_* store files.
```

The application accepts two commands: `get` and `put`. The `put` command takes
two arguments: a _key_ and a string _value_. The key will be appended to the
filename, while the value will be saved to that file.

Each plugin has its own filename prefix, e.g. `plugin-go-grpc` uses `kv_grpc_`.

Here's a full example using the `plugin-go-grpc` plugin:

```sh
# Writes to the file: kv_grpc_hello
$ ./app --grpc put hello "big wide world"

$ ./app --grpc get hello
big wide world

Written from plugin-go-grpc
```

The other plugins can be used by changing `--grpc` to either `--rpc` or `--python`.


## LICENSE

All new code and documentation, copyright (c) 2023 Michael R. Cook.

Based on the examples from https://github.com/hashicorp/go-plugin, copyright (c) HashiCorp, Inc.

SPDX-License-Identifier: MPL-2.0
