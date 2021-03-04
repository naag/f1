package rpc

import (
	"net/rpc"
)

type RpcClient struct {
	client *rpc.Client
}

func (c *RpcClient) GetScenarios() ([]string, error) {
	var resp []string

	err := c.client.Call("Plugin.GetScenarios", new(interface{}), &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *RpcClient) SetupScenario(name string) error {
	var err error
	return c.client.Call("Plugin.SetupScenario", name, err)
}

func (c *RpcClient) RunScenarioIteration(name string) error {
	var err error
	return c.client.Call("Plugin.RunScenarioIteration", name, err)
}

func (c *RpcClient) StopScenario(name string) error {
	var err error
	return c.client.Call("Plugin.StopScenario", name, err)
}
