# Bi-directional gRPC HashiCorp go-plugin Example

This example builds a simple key/counter store CLI where the mechanism
for storing and retrieving keys is pluggable. However, in this example we don't
trust the plugin to do the summation work, so we use bi-directional plugins to
call back into the main process to do the sum of two numbers.

The plugin implementation communicates over gRPC.

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
two arguments: a _key_ and a number _value_. The key will be appended to the
filename, while the number will be added to that already present in the file.

Here's a full example:

```sh
# This tells the app to use the plugin binary
$ export COUNTER_PLUGIN="./counter-go-grpc"

$ ./app put socks 2
$ ./app get socks
```

## LICENSE

All new code and documentation, copyright (c) 2023 Michael R. Cook.

Based on the examples from https://github.com/hashicorp/go-plugin, copyright (c) HashiCorp, Inc.

SPDX-License-Identifier: MPL-2.0
