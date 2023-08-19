package sbsdk

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type ProviderPlugin struct {
	Impl Provider
}

func (p *ProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ProviderRPCServer{Impl: p.Impl}, nil
}

func (p *ProviderPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ProviderRPCClient{client: c}, nil
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "Switchboard",
	MagicCookieValue: "Plugin",
}
