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

	readingsSource chan []probe.Reading
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

	dbMonitor, err := tx.LoadMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load monitor %d", id)
	}
	if dbMonitor == nil {
		return nil, errors.Errorf("no monitor with ID %d", id)
	}

	readingsChan := make(chan []probe.Reading)

	probe, err := probe.New(dbMonitor.ProbeType, dbMonitor.Probe, readingsChan)
	if err != nil {
		return nil, errors.Maskf(err, "make probe for monitor %d", id)
	}

	dbTriggers, err := tx.LoadTriggersForMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load triggers for monitor %d", id)
	}

	triggers := make([]*monitorTrigger, 0, len(dbTriggers))
	for _, dbTrigger := range dbTriggers {
		r, err := regexp.Compile(dbTrigger.Subprobes)
		if err != nil {
			// TODO(eefi): Log the problem.
			continue
		}

		t := &monitorTrigger{
			subprobes: r,
			trigger: trigger{
				Trigger: dbTrigger.Trigger,
				Env:     env,
			},
		}
		triggers = append(triggers, t)
	}

	return &monitor{
		Monitor:        dbMonitor,
		probe:          probe,
		triggers:       triggers,
		readingsSource: readingsChan,
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

func (m *monitor) process(readings []probe.Reading) {
	fmt.Printf("monitor %d got Readings %s\n", m.MonitorID, readings)
}

func (m *monitor) stop() {
	m.probe.Stop()
	close(m.readingsSource)
	<-m.stopped
}
