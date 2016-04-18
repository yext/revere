// Package daemon implements the core engine for Revere's daemon mode.
package daemon

import (
	"fmt"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
)

// Daemon represents the part of Revere that actually executes monitors and
// triggers and dispatches alerts.
type Daemon struct {
	monitors map[db.MonitorID]*monitor
	*env.Env
}

// New initializes a new Daemon. To actually make the Daemon run, call Start.
func New(env *env.Env) *Daemon {
	return &Daemon{monitors: make(map[db.MonitorID]*monitor), Env: env}
}

// Start starts running a Daemon.
func (d *Daemon) Start() {
	monitorIDs, err := d.DB.LoadAllMonitorIDs()
	if err != nil {
		// TODO(eefi): Change when implementing monitor reloading.
		panic(fmt.Sprintf("start daemon: load monitor IDs: %v", err))
	}

	for _, id := range monitorIDs {
		monitor, err := newMonitor(id, d.Env)
		if err != nil {
			// TODO(eefi): Change when implementing monitor reloading.
			panic(fmt.Sprintf("start daemon: load monitor %d: %v", id, err))
		}

		d.monitors[id] = monitor
		monitor.start()
	}
}

// Stop gracefully stops a Daemon. It tries to allow any in-progress delivery of
// alerts to finish before returning.
func (d *Daemon) Stop() {
	for _, m := range d.monitors {
		m.stop()
	}
}
