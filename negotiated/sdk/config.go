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
	// This isn't required as were using VersionedPlugins.
	// ProtocolVersion: 1,

	// Once set, these magic cookie values should NEVER be changed.
	MagicCookieKey:   "NEGOTIATED_PLUGIN", // a unique key for your application
	MagicCookieValue: "negotiated",        // a random 256-bit value is recommended
}