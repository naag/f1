package plugin

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/pkg/api"
)

type ScenarioPlugin struct {
	scenarios map[string]*scenario
}

type scenario struct {
	setupFn    api.SetupFn
	runFn      api.RunFn
	teardownFn api.TeardownFn
	t          *api.T
}

func NewServer() *ScenarioPlugin {
	p := &ScenarioPlugin{}
	p.scenarios = make(map[string]*scenario)
	return p
}

func (p *ScenarioPlugin) WithScenario(name string, setupFn api.SetupFn) *ScenarioPlugin {
	p.scenarios[name] = &scenario{
		setupFn: setupFn,
	}
	return p
}

func (p *ScenarioPlugin) Serve() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	pluginMap := map[string]plugin.Plugin{
		"scenarioPlugin": &api.ScenarioPlugin{Impl: p},
	}

	handshakeConfig := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}
