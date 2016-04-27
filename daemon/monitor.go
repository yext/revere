package daemon

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
	"github.com/yext/revere/probe"
)

type monitor struct {
	*db.Monitor

	probe    probe.Probe
	triggers []monitorTrigger

	readingsSource chan []probe.Reading
	stopped        chan struct{}

	*env.Env
}

type monitorTrigger struct {
	subprobes *regexp.Regexp
	*trigger
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

	dbMonitorTriggers, err := tx.LoadTriggersForMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load triggers for monitor %d", id)
	}

	monitorTriggers := make([]monitorTrigger, 0, len(dbMonitorTriggers))
	for _, dbMonitorTrigger := range dbMonitorTriggers {
		monitorTrigger, err := newMonitorTrigger(dbMonitorTrigger)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"monitor": id,
				"trigger": dbMonitorTrigger.TriggerID,
			}).Error("Could not load monitor trigger. Discarding.")
			continue
		}
		monitorTriggers = append(monitorTriggers, *monitorTrigger)
	}

	return &monitor{
		Monitor:        dbMonitor,
		probe:          probe,
		triggers:       monitorTriggers,
		readingsSource: readingsChan,
		stopped:        make(chan struct{}),
		Env:            env,
	}, nil
}

func newMonitorTrigger(dbMonitorTrigger db.MonitorTrigger) (*monitorTrigger, error) {
	subprobesRegexp, err := regexp.Compile(dbMonitorTrigger.Subprobes)
	if err != nil {
		return nil, errors.Maskf(err, "compile regexp")
	}

	trigger, err := newTrigger(dbMonitorTrigger.Trigger)
	if err != nil {
		return nil, errors.Maskf(err, "make trigger")
	}

	return &monitorTrigger{
		subprobes: subprobesRegexp,
		trigger:   trigger,
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
