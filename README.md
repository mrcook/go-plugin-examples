# Hashicorp go-plugin examples

This repository contains the example Go code from HashiCorps `go-plugin`, which
have been refactored and documented with the hope of making them a little easier
to follow and understand when using in your own plugin based applications.

* `basic`: a simple example with communication over net/rpc
* `bidirectional`: a gRPC example with two-way communication between host <-> plugin
* `gprc`: an example with communication over gRPC, including Go and Python plugin examples
* `negotiated`: an example handling different versions of the same plugin: one using net/rpc, the other gRPC.

## LICENSE

All new code and documentation, copyright (c) 2023 Michael R. Cook.

Based on the examples from https://github.com/hashicorp/go-plugin, copyright (c) HashiCorp, Inc.

SPDX-License-Identifier: MPL-2.0
