package plugin

import (
	"os"
	"path"

	"github.com/form3tech-oss/f1/pkg/f1/testing"
	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/pkg/api"
	"github.com/naag/f1-api/pkg/client"
)

func RegisterPlugin(p api.ScenarioPluginInterface, pluginName string) error {
	scenarios, err := p.GetScenarios()
	if err != nil {
		return err
	}

	for _, scenarioName := range scenarios {
		remoteName := scenarioName
		localName := pluginName + "/" + remoteName

		setupFn := func(t *testing.T) (testing.RunFn, testing.TeardownFn) {
			err := p.SetupScenario(remoteName)
			if err != nil {
				t.Fail()
			}

			runFn := func(t *testing.T) {
				err := p.RunScenarioIteration(remoteName)
				if err != nil {
					t.Fail()
				}
			}

			teardownFn := func(t *testing.T) {
				err := p.StopScenario(remoteName)
				if err != nil {
					t.Fail()
				}
			}

			return runFn, teardownFn
		}

		testing.Add(localName, setupFn)
	}

	return nil
}

func DiscoverPlugins() ([]string, error) {
	return plugin.Discover("*-plugin", "~/.f1/plugins")
}

func LaunchAll(verbose bool) (func(), error) {
	var clients []*plugin.Client

	paths, err := plugin.Discover("*-plugin", pluginDir())
	if err != nil {
		return nil, err
	}

	for _, pluginPath := range paths {
		c, err := client.NewClient().
			WithVerboseLogging(verbose).
			Build(pluginPath).
			Connect()

		if err != nil {
			return nil, err
		}

		clients = append(clients, c.Client)

		err = RegisterPlugin(c.Impl, path.Base(pluginPath))
		if err != nil {
			return nil, err
		}
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
