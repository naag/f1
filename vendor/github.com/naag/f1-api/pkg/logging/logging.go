package logging

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

func GetServerLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       "pluginServer",
		Output:     os.Stderr,
		Level:      hclog.Info,
		JSONFormat: true,
	})
}

func GetClientLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       "pluginClient",
		Output:     os.Stderr,
		Level:      hclog.Info,
		JSONFormat: false,
	})
}
