package marathon

// Application represents the object for an application in marathon
type Application struct {
	ID              string             `json:"id,omitempty"`
	Instances       *int               `json:"instances,omitempty"`
	Tasks           []*Task            `json:"tasks,omitempty"`
	Ports           []int              `json:"ports"`
	PortDefinitions *[]PortDefinition  `json:"portDefinitions,omitempty"`
	Labels          *map[string]string `json:"labels,omitempty"`
}

// PortDefinition represents a port that should be considered part of
// a resource. Port definitions are necessary when you are using HOST
// networking
type PortDefinition struct {
	Port     *int               `json:"port,omitempty"`
	Protocol string             `json:"protocol,omitempty"`
	Name     string             `json:"name,omitempty"`
	Labels   *map[string]string `json:"labels,omitempty"`
}

// Task represents the definition for a marathon task
type Task struct {
	ID           string       `json:"id"`
	AppID        string       `json:"appId"`
	Host         string       `json:"host"`
	Ports        []int        `json:"ports"`
	ServicePorts []int        `json:"servicePorts"`
	SlaveID      string       `json:"slaveId"`
	StagedAt     string       `json:"stagedAt"`
	StartedAt    string       `json:"startedAt"`
	State        string       `json:"state"`
	IPAddresses  []*IPAddress `json:"ipAddresses"`
	Version      string       `json:"version"`
}

// IPAddress represents a task's IP address and protocol.
type IPAddress struct {
	IPAddress string `json:"ipAddress"`
	Protocol  string `json:"protocol"`
}
