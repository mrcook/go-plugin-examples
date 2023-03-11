# Basic HashiCorp go-plugin Example

A simple example of a `go-plugin` plugin that communicates over net/rpc.

This example contains the host application code, the example plugin
implementation, along with a plugin SDK that is used by the host application
and plugin for communicating with the `Greeter` plugin type.

## hello_plugin_example

An example plugin.

This contains `HelloGreeterPlugin`, a real implementation of the `Greeter`
plugin type, and a `main` function entry point.

## sdk

The Plugin Software Development Kit.

This package (which could be released as a standalone library) is used by
plugin authors for creating `Greeter` plugins. It provides the `GreeterPlugin`,
which is an implementation of the `plugin.Plugin` that `go-plugin` uses to
allow communication between the host application and plugins.


## Usage

A `Makefile` is provide for ease of use:

    make

This will compile both the host application and plugin, and then run the application.

The output shows the server logs including the logs called by the plugin. Along
with this the application will print the greeting that the plugin returned.

Additional make commands:

    make build
    make run
    make clean


## LICENSE

All new code and documentation, copyright (c) 2023 Michael R. Cook.

Based on the examples from https://github.com/hashicorp/go-plugin, copyright (c) HashiCorp, Inc.

SPDX-License-Identifier: MPL-2.0
