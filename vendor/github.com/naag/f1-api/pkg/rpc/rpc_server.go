package rpc

import (
	"github.com/naag/f1-api/pkg/api"
)

type RpcServer struct {
	Impl api.ScenarioPluginInterface
}

func (s *RpcServer) GetScenarios(args interface{}, resp *[]string) error {
	var err error
	*resp, err = s.Impl.GetScenarios()
	return err
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
