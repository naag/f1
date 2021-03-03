package api

import "github.com/hashicorp/go-plugin"

type RpcServer struct {
	Impl ScenarioPluginInterface
}

func (s *RpcServer) GetScenarios(args interface{}, resp *[]string) error {
	*resp = s.Impl.GetScenarios()
	return nil
}

func (s *RpcServer) SetupScenario(name string, resp *[]string) error {
	return s.Impl.SetupScenario(name)
}

func (s *RpcServer) RunScenarioIteration(name string, resp *[]string) error {
	return s.Impl.RunScenarioIteration(name)
}

func (s *RpcServer) StopScenario(name string, resp *[]string) error {
	return s.Impl.StopScenario(name)
}

func (p *ScenarioPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RpcServer{Impl: p.Impl}, nil
}
