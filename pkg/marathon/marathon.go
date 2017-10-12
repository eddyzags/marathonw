package marathon

import "fmt"

// Marathon is an interface to provide a http client for marathon framework
type Marathon interface {
	Ping() error
	Applications(label string) ([]*Application, error)
	Tasks(appID string) ([]*Task, error)
}

type marathon struct {
	config *Config
}

// Config represents the marathon client configuration object
type Config struct {
	HTTPBasicAuthUser     string
	HTTPBasicAuthPassword string
	DCOSToken             string
	URI                   string
}

// NewClient instantiates a new marathon client
func NewClient(config *Config) Marathon {
	return &marathon{
		config: config,
	}
}

// Applications returns a set of applications according to a label
func (m *marathon) Applications(label string) ([]*Application, error) {
	apps := []*Application{}

	path := fmt.Sprintf("/v2/apps?label=%s", label)

	if err := m.apiCall("GET", path, nil, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

// Tasks returns a specific application's set of tasks
func (m *marathon) Tasks(appID string) ([]*Task, error) {
	tasks := []*Task{}

	if err := m.apiCall("GET", "/v2/tasks", nil, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Ping returns an error if the marathon framework is unreachable
func (m *marathon) Ping() error {
	return m.apiCall("GET", "/ping", nil, nil)
}

// URI returns the marathon uri
func (m *marathon) URI() string {
	return m.config.URI
}
