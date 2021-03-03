package api

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type RpcClient struct {
	client *rpc.Client
}

func (g *RpcClient) GetScenarios() []string {
	var resp []string
	err := g.client.Call("Plugin.GetScenarios", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (g *RpcClient) SetupScenario(name string) error {
	var err error

	clientErr := g.client.Call("Plugin.SetupScenario", name, err)
	if clientErr != nil {
		panic(clientErr)
	}

	return err
}

func (g *RpcClient) RunScenarioIteration(name string) error {
	var err error
	g.client.Call("Plugin.RunScenarioIteration", name, err)
	return err
}

func (g *RpcClient) StopScenario(name string) error {
	var err error

	clientErr := g.client.Call("Plugin.StopScenario", name, err)
	if clientErr != nil {
		panic(clientErr)
	}

	return err
}

func (ScenarioPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RpcClient{client: c}, nil
}
