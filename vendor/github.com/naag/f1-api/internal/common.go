package internal

import "github.com/hashicorp/go-plugin"

const (
	PluginName = "scenarioPlugin"
)

var (
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}
)
