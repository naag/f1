package client

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/internal"
	"github.com/naag/f1-api/pkg/api"
	"github.com/naag/f1-api/pkg/logging"
	"github.com/naag/f1-api/pkg/rpc"
)

type PluginClient struct {
	Client  *plugin.Client
	Impl    api.ScenarioPluginInterface
	verbose bool
}

func NewClient() *PluginClient {
	return &PluginClient{
		verbose: false,
	}
}

func (p *PluginClient) WithVerboseLogging(verbose bool) *PluginClient {
	p.verbose = verbose
	return p
}

func (p *PluginClient) Build(pluginPath string) *PluginClient {
	logger := logging.GetClientLogger()

	pluginMap := map[string]plugin.Plugin{
		internal.PluginName: &rpc.Plugin{},
	}

	stdoutCopier := newCopier(os.Stdout)
	stderrCopier := newCopier(os.Stderr)

	p.Client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: internal.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             p.getCmd(pluginPath),
		Logger:          logger,
		SyncStdout:      stdoutCopier,
		SyncStderr:      stderrCopier,
	})

	stdoutCopier.Sync()
	stderrCopier.Sync()

	// 	defer func() {
	// 		_ = stdoutCopier.Close()
	// 		_ = stderrCopier.Close()
	// 	}()

	return p
}

func (p *PluginClient) Connect() (*PluginClient, error) {
	rpcClient, err := p.Client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(internal.PluginName)
	if err != nil {
		return nil, err
	}

	p.Impl = raw.(api.ScenarioPluginInterface)

	return p, nil
}

func (p *PluginClient) getCmd(pluginPath string) *exec.Cmd {
	if p.verbose {
		return exec.Command(pluginPath, "-v")
	} else {
		return exec.Command(pluginPath)
	}
}

type copier struct {
	data       []byte
	sniff      []byte
	target     io.Writer
	mu         sync.Mutex
	closing    bool
	doneChan   chan error
	exitString string
}

const exitStr = "\n{plugin:terminated}\r\x1b[2K"

func newCopier(dst io.Writer) *copier {
	return &copier{
		target:     dst,
		doneChan:   make(chan error),
		exitString: exitStr,
	}
}

func (c *copier) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = append(c.data, b...)
	c.sniff = append(c.sniff, b...)
	if len(c.sniff) > len(c.exitString) {
		c.sniff = c.sniff[len(c.sniff)-len(c.exitString):]
	}
	return len(b), nil
}

func (c *copier) Close() error {
	c.mu.Lock()
	c.closing = true
	c.mu.Unlock()
	return <-c.doneChan
}

func (c *copier) Sync() {

	lastDataAt := time.Now()

	go func() {
		for {
			if func() bool {

				c.mu.Lock()
				defer c.mu.Unlock()

				if len(c.data) == 0 {
					if c.closing && (string(c.sniff) == c.exitString || time.Since(lastDataAt) > time.Second*10) {
						c.doneChan <- nil
						return true
					}
					time.Sleep(time.Millisecond * 10)
					return false
				}

				lastDataAt = time.Now()

				count, err := c.target.Write(c.data)
				if err != nil {
					c.doneChan <- err
					return true
				}
				c.data = c.data[count:]
				return false
			}() {
				break
			}
		}
	}()
}
