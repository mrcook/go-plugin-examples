// Copyright (c) Michael R. Cook.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import "github.com/hashicorp/go-plugin"

// HandshakeConfig is used to perform a basic handshake between a plugin and
// host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var HandshakeConfig = plugin.HandshakeConfig{
	// In general it is better to omit ProtocolVersion and explicitly set
	// VersionedPlugins in Client/Server configurations.
	// For simplicity, ProtocolVersion is used in this example.
	ProtocolVersion: 1,

	// Once set, these magic cookie values should NEVER be changed.
	MagicCookieKey:   "GRPC_PLUGIN", // a unique key for your application
	MagicCookieValue: "grpc",        // a random 256-bit value
}
