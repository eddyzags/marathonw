package marathonw

import (
	"google.golang.org/grpc/naming"
)

// Resolver is an object for managing watcher
type Resolver struct {
	*Config
}

// NewResolver instantiates a new naming resolver for a grpc client
func NewResolver(c *Config) naming.Resolver {
	return &Resolver{c}
}

// Resolve creates a watcher in order to retrieve the service name
// addresses for the load balancing process
func (r *Resolver) Resolve(name string) (naming.Watcher, error) {
	r.serviceName = name
	return NewWatcher(r.Config)
}
