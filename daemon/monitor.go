package daemon

import (
	"fmt"
	"regexp"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
)

type monitor struct {
	*db.Monitor

	triggers []*monitorTrigger

	*env.Env
}

type monitorTrigger struct {
	subprobes *regexp.Regexp
	trigger
}

func newMonitor(id db.MonitorID, env *env.Env) (*monitor, error) {
	tx, err := env.DB.Beginx()
	if err != nil {
		return nil, errors.Mask(err)
	}
	defer tx.Rollback()

	m, err := tx.LoadMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load monitor %d", id)
	}
	if m == nil {
		return nil, errors.Errorf("no monitor with ID %d", id)
	}

	dbMTs, err := tx.LoadTriggersForMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load triggers for monitor %d", id)
	}

	mts := make([]*monitorTrigger, 0, len(dbMTs))
	for _, dbMT := range dbMTs {
		r, err := regexp.Compile(dbMT.Subprobes)
		if err != nil {
			// TODO(eefi): Log the problem.
			continue
		}

		mt := &monitorTrigger{
			subprobes: r,
			trigger: trigger{
				Trigger: dbMT.Trigger,
				Env:     env,
			},
		}
		mts = append(mts, mt)
	}

	return &monitor{Monitor: m, triggers: mts, Env: env}, nil
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
