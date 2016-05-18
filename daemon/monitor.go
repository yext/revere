package daemon

import (
	"regexp"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/state"
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

	dbLabelTriggers, err := tx.LoadLabelTriggersForMonitor(id)
	if err != nil {
		return nil, errors.Maskf(err, "load label triggers for monitor %d", id)
	}

	monitorTriggers := make(
		[]monitorTrigger, 0, len(dbMonitorTriggers)+len(dbLabelTriggers))
	for _, dbMonitorTrigger := range dbMonitorTriggers {
		monitorTrigger, err := newMonitorTrigger(
			dbMonitorTrigger.Subprobes, dbMonitorTrigger.Trigger)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"monitor": id,
				"trigger": dbMonitorTrigger.TriggerID,
			}).Error("Could not load monitor trigger. Discarding.")
			continue
		}
		monitorTriggers = append(monitorTriggers, *monitorTrigger)
	}
	for _, dbLabelTrigger := range dbLabelTriggers {
		monitorTrigger, err := newMonitorTrigger(
			dbLabelTrigger.Subprobes, dbLabelTrigger.Trigger)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"monitor": id,
				"label":   dbLabelTrigger.LabelID,
				"trigger": dbLabelTrigger.TriggerID,
			}).Error("Could not load label trigger. Discarding.")
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

func newMonitorTrigger(subprobes string, dbTrigger *db.Trigger) (*monitorTrigger, error) {
	subprobesRegexp, err := regexp.Compile(subprobes)
	if err != nil {
		return nil, errors.Maskf(err, "compile regexp")
	}

	triggerTemplate, err := newTriggerTemplate(dbTrigger)
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
	var silences []silence
	if m.shouldLoadSilences(readings) {
		silences = m.loadActiveSilences()
	}

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

		isSilenced := false
		for _, silence := range silences {
			if silence.silences(subprobe) {
				isSilenced = true
				break
			}
		}

		subprobe.process(r, isSilenced)
	}
}

// shouldLoadSilences returns whether processing the given set of readings
// against the current state of this monitor's subprobes might require checking
// for silences.
//
// In the common case of all subprobes currently normal and all incoming
// readings reading normal, no alerts will need to be sent, so it doesn't matter
// whether there are any active silences. We can avoid a DB round trip when this
// is the case.
func (m *monitor) shouldLoadSilences(readings []probe.Reading) bool {
	for _, s := range m.subprobes {
		if s.state != state.Normal {
			return true
		}
	}
	for _, r := range readings {
		if r.State != state.Normal {
			return true
		}
	}
	return false
}

func (m *monitor) loadActiveSilences() []silence {
	dbSilences, err := m.DB.LoadActiveSilencesForMonitor(m.id)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"monitor": m.id,
		}).Error("Could not load active silences. Proceeding without silencing.")
		return nil
	}

	if len(dbSilences) == 0 {
		return nil
	}

	silences := make([]silence, 0, len(dbSilences))
	for _, dbSilence := range dbSilences {
		s, err := newSilence(dbSilence)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"monitor": m.id,
				"silence": dbSilence.SilenceID,
			}).Error("Could not load silence. Ignoring.")
			continue
		}
		silences = append(silences, s)
	}
	return silences
}

func (m *monitor) stop() {
	m.stopper.Do(func() {
		m.probe.Stop()
		close(m.readingsSource)
		<-m.stopped
	})
}
