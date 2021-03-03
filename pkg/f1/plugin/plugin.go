package plugin

import (
	"os"

	"github.com/form3tech-oss/f1/pkg/f1/testing"
	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/pkg/api"
	f1plugin "github.com/naag/f1-api/pkg/plugin"
)

func RegisterPlugin(p api.ScenarioPluginInterface) {
	for _, scenarioName := range p.GetScenarios() {
		copyScenario := scenarioName
		setupFn := func(t *testing.T) (testing.RunFn, testing.TeardownFn) {
			p.SetupScenario(copyScenario)

			runFn := func(t *testing.T) {
				p.RunScenarioIteration(copyScenario)
			}

			teardownFn := func(t *testing.T) {
				p.StopScenario(copyScenario)
			}

			return runFn, teardownFn
		}

		testing.Add(copyScenario, setupFn)
	}
}

func DiscoverPlugins() ([]string, error) {
	return plugin.Discover("*-plugin", "~/.f1/plugins")
}

func LaunchAll() (func(), error) {
	var clients []*plugin.Client

	paths, err := plugin.Discover("*-plugin", pluginDir())
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		c, p, err := f1plugin.NewClient(path)
		if err != nil {
			return nil, err
		}

		clients = append(clients, c)
		RegisterPlugin(p)
	}

	fn := func() {
		for _, c := range clients {
			c.Kill()
		}
	}

	return fn, nil
}

func pluginDir() string {
	return os.ExpandEnv("${HOME}/.f1/plugins")
}
