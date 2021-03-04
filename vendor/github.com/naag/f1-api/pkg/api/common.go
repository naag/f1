package api

type ScenarioPluginInterface interface {
	GetScenarios() ([]string, error)
	SetupScenario(name string) error        // Setup pool of go workers and run SetupFn
	RunScenarioIteration(name string) error // Run RunFn inside of go worker
	StopScenario(name string) error
}
