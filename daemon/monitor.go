package daemon

import (
	"regexp"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
	"github.com/yext/revere/probe"
)

type monitor struct {
	id          db.MonitorID
	name        string
	description string
	response    string
	version     int32

	probe    probe.Probe
	triggers []monitorTrigger

	subprobes map[string]*subprobe

	readingsSource chan []probe.Reading
	stopper        sync.Once
	stopped        chan struct{}

	*env.Env
}

type monitorTrigger struct {
	subprobes *regexp.Regexp
	*triggerTemplate
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

	monitor := &monitor{
		id:             id,
		name:           dbMonitor.Name,
		description:    dbMonitor.Description,
		response:       dbMonitor.Response,
		version:        dbMonitor.Version,
		probe:          probe,
		triggers:       monitorTriggers,
		subprobes:      make(map[string]*subprobe),
		readingsSource: readingsChan,
		stopped:        make(chan struct{}),
		Env:            env,
	}

	dbSubprobeStatuses, err := tx.LoadSubprobeStatusesForMonitor(id)
	if err != nil {
		// It's possible to still generate alerts even with brokenness
		// in saving state to the DB, so log and stumble on.
		log.WithError(err).WithFields(log.Fields{
			"monitor": id,
		}).Error("Could not load subprobe statuses. Alerts might have inaccurate historical data.")
		return monitor, nil
	}

	for name, status := range dbSubprobeStatuses {
		monitor.subprobes[name] = newSubprobe(name, status, monitor)
	}

	return monitor, nil
}

func newMonitorTrigger(dbMonitorTrigger db.MonitorTrigger) (*monitorTrigger, error) {
	subprobesRegexp, err := regexp.Compile(dbMonitorTrigger.Subprobes)
	if err != nil {
		return nil, errors.Maskf(err, "compile regexp")
	}

	triggerTemplate, err := newTriggerTemplate(dbMonitorTrigger.Trigger)
	if err != nil {
		return nil, errors.Maskf(err, "make trigger")
	}

	return &monitorTrigger{
		subprobes:       subprobesRegexp,
		triggerTemplate: triggerTemplate,
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
	for _, r := range readings {
		subprobe := m.subprobes[r.Subprobe]
		if subprobe == nil {
			var err error
			subprobe, err = createSubprobe(m, r)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"monitor":  m.id,
					"subprobe": r.Subprobe,
					"state":    r.State,
					"recorded": r.Recorded,
				}).Error("Could not create subprobe. Discarding reading.")
				continue
			}

			m.subprobes[subprobe.name] = subprobe
		}

		subprobe.process(r)
	}
}

func (m *monitor) stop() {
	m.stopper.Do(func() {
		m.probe.Stop()
		close(m.readingsSource)
		<-m.stopped
	})
}
