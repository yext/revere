package daemon

import (
	"fmt"
	"regexp"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
	"github.com/yext/revere/probe"
)

type monitor struct {
	*db.Monitor

	probe          probe.Probe
	readingsSource <-chan *probe.Readings

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

	c := make(chan *probe.Readings)

	p, err := probe.New(m.ProbeType, m.Probe, c)
	if err != nil {
		return nil, errors.Maskf(err, "make probe for monitor %d", id)
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

	return &monitor{
		Monitor:        m,
		probe:          p,
		readingsSource: c,
		triggers:       mts,
		Env:            env,
	}, nil
}

func (m *monitor) start() {
	// TODO(eefi): Implement.
	fmt.Printf("starting monitor %d\n", m.MonitorID)

	m.probe.Start()
}

func (m *monitor) stop() {
	// TODO(eefi): Implement.
	fmt.Printf("stopping monitor %d\n", m.MonitorID)

	m.probe.Stop()
}
