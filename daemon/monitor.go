package daemon

import (
	"fmt"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
)

type monitor struct {
	*db.Monitor
	*env.Env
}

func newMonitor(id db.MonitorID, env *env.Env) (*monitor, error) {
	m, err := env.DB.LoadMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load monitor %d", id)
	}

	if m == nil {
		return nil, errors.Errorf("no monitor with ID %d", id)
	}

	return &monitor{Monitor: m, Env: env}, nil
}

// Start starts running a monitor.
func (m *monitor) Start() {
	// TODO(eefi): Implement.
	fmt.Printf("starting monitor %d\n", m.MonitorID)
}

// Stop gracefully stops a monitor.
func (m *monitor) Stop() {
	// TODO(eefi): Implement.
	fmt.Printf("stopping monitor %d\n", m.MonitorID)
}
