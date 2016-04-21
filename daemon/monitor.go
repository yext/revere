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

	probe    probe.Probe
	triggers []*monitorTrigger

	readingsSource chan *probe.Readings
	stopped        chan struct{}

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
		triggers:       mts,
		readingsSource: c,
		stopped:        make(chan struct{}),
		Env:            env,
	}, nil
}

func (m *monitor) start() {
	m.probe.Start()
	go func() {
		defer close(m.stopped)
		for {
			r, ok := <-m.readingsSource
			if !ok {
				return
			}
			m.process(r)
		}
	}()
}

func (m *monitor) process(readings *probe.Readings) {
	fmt.Printf("monitor %d got Readings %s\n", m.MonitorID, readings)
}

func (m *monitor) stop() {
	m.probe.Stop()
	close(m.readingsSource)
	<-m.stopped
}
