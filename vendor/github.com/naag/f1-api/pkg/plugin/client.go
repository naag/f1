package plugin

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/pkg/api"
)

func NewClient(pluginPath string) (*plugin.Client, api.ScenarioPluginInterface, error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	pluginMap := map[string]plugin.Plugin{
		"scenarioPlugin": &api.ScenarioPlugin{},
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}

	raw, err := rpcClient.Dispense("scenarioPlugin")
	if err != nil {
		return nil, nil, err
	}

	return client, raw.(api.ScenarioPluginInterface), nil
}
