package rpc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/pkg/api"
)

type Plugin struct {
	Impl api.ScenarioPluginInterface
}

func (Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{client: c}, nil
}

func (p *Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RpcServer{Impl: p.Impl}, nil
}
