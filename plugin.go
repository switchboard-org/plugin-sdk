package main

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
