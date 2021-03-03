package plugin

import (
	"errors"

	"github.com/naag/f1-api/pkg/api"
)

func (g *ScenarioPlugin) GetScenarios() []string {
	var result []string
	for name := range g.scenarios {
		result = append(result, name)
	}
	return result
}

func (g *ScenarioPlugin) SetupScenario(name string) error {
	s := g.getScenarioByName(name)
	s.t = newT(name)
	s.runFn, s.teardownFn = s.setupFn(s.t)

	if s.t.HasFailed() {
		return errors.New("setup scenario failed")
	}

	return nil
}

func (g *ScenarioPlugin) RunScenarioIteration(name string) error {
	s := g.getScenarioByName(name)
	s.t = newT(name)

	s.runFn(s.t)

	if s.t.HasFailed() {
		return errors.New("iteration failed")
	}

	return nil
}

func (g *ScenarioPlugin) StopScenario(name string) error {
	s := g.getScenarioByName(name)
	s.teardownFn(s.t)

	if s.t.HasFailed() {
		return errors.New("stop scenario failed")
	}

	return nil
}

func (g *ScenarioPlugin) getScenarioByName(name string) *scenario {
	return g.scenarios[name]
}

func newT(name string) *api.T {
	return api.NewT(api.LoadEnvironment(), "0", "0", name)
}
