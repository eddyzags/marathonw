package marathonw

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eddyzags/marathonw/pkg/marathon"

	"google.golang.org/grpc/naming"
)

// Watcher watches the updates on a specific service target
type Watcher struct {
	serviceName    string
	marathonClient marathon.Marathon
	addrs          []string
	appID          string
	portIndex      int
	sync.Mutex
}

// Config is the watcher configuration object
type Config struct {
	MarathonURI  string
	serviceName  string
	DCOSUser     string
	DCOSPassword string
	DCOSToken    string
}

// NewWatcher instantiates a new marathon watcher
func NewWatcher(config *Config) (*Watcher, error) {
	c := &marathon.Config{
		HTTPBasicAuthUser:     config.DCOSUser,
		HTTPBasicAuthPassword: config.DCOSPassword,
		DCOSToken:             config.DCOSToken,
	}

	if config.MarathonURI == "" {
		config.MarathonURI = "marathon.mesos:8080"
	}

	marathonClient := marathon.NewClient(c)
	if err := marathonClient.Ping(); err != nil {
		return nil, err
	}

	apps, err := marathonClient.Applications(config.serviceName)
	if err != nil {
		return nil, err
	}

	// Only one application must be found
	if len(apps) != 1 {
		return nil, ErrServiceNameCollision
	}

	portIndex, err := findPortIndex(config.serviceName, apps[0])
	if err != nil {
		return nil, err
	}

	return &Watcher{
		serviceName:    config.serviceName,
		marathonClient: marathonClient,
		portIndex:      portIndex,
	}, nil
}

// Next calls marathon framework in order to detect a task update for a
// specific application and port index. It blocks until an update or errors happens
func (w *Watcher) Next() ([]*naming.Update, error) {
	ups := w.updates()
	if len(w.addrs) <= 0 {
		return ups, nil
	}

	ticker := time.NewTicker(2 * time.Second)

	for _ = range ticker.C {
		ups := w.updates()
		if len(ups) > 0 {
			return ups, nil
		}
	}

	return nil, nil
}

// Close closes the watcher's poll
func (w *Watcher) Close() {}

func (w *Watcher) updates() []*naming.Update {
	latest := w.poll()

	var ups []*naming.Update

	// Detecting new addrs
	for _, addr := range latest {
		if !contains(w.addrs, addr) {
			ups = append(ups, &naming.Update{Op: naming.Add, Addr: addr})
		}
	}

	// Detecting old addrs
	for _, addr := range w.addrs {
		if !contains(latest, addr) {
			ups = append(ups, &naming.Update{Op: naming.Delete, Addr: addr})
		}
	}

	w.addrs = latest
	return ups
}

func (w *Watcher) poll() []string {
	var addrs []string

	tasks, err := w.marathonClient.Tasks(w.appID)
	if err != nil {
		//TODO(eddyzags): Add debug log
	}

	for _, task := range tasks {
		addrs = append(addrs,
			task.Host+":"+strconv.FormatInt(int64(task.Ports[w.portIndex]), 10))
	}

	return addrs
}

func findPortIndex(serviceName string, app *marathon.Application) (int, error) {
	for k, v := range *app.Labels {
		if v == serviceName && strings.HasPrefix(k, "MARATHONW_") && strings.HasSuffix(k, "_NAME") {
			return strconv.Atoi(strings.Trim(strings.Trim(k, "MARATHONW_"), "_NAME"))
		}
	}

	return -1, ErrServiceLabelNotFound
}

func contains(array []string, s string) bool {
	for _, a := range array {
		if a == s {
			return true
		}
	}
	return false
}
