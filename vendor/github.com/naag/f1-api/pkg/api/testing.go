package api

import (
	"os"
	"runtime"
	"strings"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type T struct {
	// Identifier of the user for the test
	VirtualUser string
	// Iteration number, "setup" or "teardown"
	Iteration string
	// Logger with user and iteration tags
	Log         *log.Logger
	failed      int64
	Require     *require.Assertions
	Environment map[string]string
	Scenario    string
}

// SetupFn performs any setup required to run a scenario.
// It returns a RunFn to be invoked for every iteration of the tests
// and a TeardownFn to clear down any resources after all iterations complete
type SetupFn func(t *T) (RunFn, TeardownFn)

// RunFn performs a single iteration of the test. It my be used for asserting
// results or failing the test.
type RunFn func(t *T)

// TeardownFn clears down any resources from a test run after all iterations complete.
type TeardownFn RunFn

func NewT(env map[string]string, vu, iter string, scenarioName string) *T {
	t := &T{
		VirtualUser: vu,
		Iteration:   iter,
		Log:         log.WithField("u", vu).WithField("i", iter).WithField("scenario", scenarioName).Logger,
		Environment: env,
		Scenario:    scenarioName,
	}
	t.Require = require.New(t)
	return t
}

func (t *T) Errorf(format string, args ...interface{}) {
	atomic.StoreInt64(&t.failed, int64(1))
	t.Log.Errorf(format, args...)
}

func (t *T) FailNow() {
	atomic.StoreInt64(&t.failed, int64(1))
	t.Log.Errorf("test failed and stopped")
	runtime.Goexit()
}

func (t *T) Fail() {
	atomic.StoreInt64(&t.failed, int64(1))
	t.Log.Errorf("test failed")
}

func (t *T) FailWithError(err error) {
	atomic.StoreInt64(&t.failed, int64(1))
	t.Log.WithError(err).Errorf("test failed due to %s", err.Error())
}

func (t *T) HasFailed() bool {
	return atomic.LoadInt64(&t.failed) == int64(1)
}

func LoadEnvironment() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		keyAndValue := strings.SplitN(e, "=", 2)
		if len(keyAndValue) < 2 {
			log.Warnf("Malformed environment variable was not loaded: %s", e)
			continue
		}
		env[keyAndValue[0]] = keyAndValue[1]
	}
	return env
}
