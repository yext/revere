package daemon

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/state"
)

type subprobe struct {
	id      db.SubprobeID
	monitor *monitor
	name    string

	lastReading  time.Time
	state        state.State
	enteredState time.Time
	lastNormal   time.Time

	saveNextReading bool

	triggerSets map[db.TargetType]*sameTypeTriggerSet

	*env.Env
}

func newSubprobe(name string, status db.SubprobeStatus, monitor *monitor) *subprobe {
	return &subprobe{
		id:              status.SubprobeID,
		monitor:         monitor,
		name:            name,
		lastReading:     status.Recorded,
		state:           status.State,
		enteredState:    status.EnteredState,
		lastNormal:      status.LastNormal,
		saveNextReading: false,
		triggerSets:     newSubprobeTriggerSets(monitor, name),
		Env:             monitor.Env,
	}
}

// createSubprobe creates a new subprobe in the database based on receiving a
// reading for a previously unknown subprobe.
func createSubprobe(monitor *monitor, reading probe.Reading) (*subprobe, error) {
	s := &subprobe{
		monitor: monitor,
		name:    reading.Subprobe,

		lastReading:  reading.Recorded,
		state:        reading.State,
		enteredState: reading.Recorded,

		// A bit of a lie if state != Normal, but it's the best we have.
		lastNormal: reading.Recorded,

		// Make sure the first reading is saved.
		saveNextReading: true,

		triggerSets: newSubprobeTriggerSets(monitor, reading.Subprobe),

		Env: monitor.Env,
	}

	err := s.DB.Tx(func(tx *db.Tx) error {
		var err error

		s.id, err = tx.InsertSubprobe(monitor.id, s.name)
		if err != nil {
			return errors.Maskf(err, "insert subprobe")
		}

		status := s.dbStatus()
		err = tx.InsertSubprobeStatus(status)
		if err != nil {
			return errors.Maskf(err, "insert subprobe status")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Mask(err)
	}

	return s, nil
}

// newSubprobeTriggerSets filters monitor's triggers down to a map appropriate
// for the triggerSets field of a subprobe with the given name.
func newSubprobeTriggerSets(monitor *monitor, name string) map[db.TargetType]*sameTypeTriggerSet {
	triggerSets := make(map[db.TargetType]*sameTypeTriggerSet)
	for _, monitorTrigger := range monitor.triggers {
		if monitorTrigger.subprobes.MatchString(name) {
			trigger := monitorTrigger.trigger
			targetType := trigger.target.Type().ID()

			triggerSet := triggerSets[targetType]
			if triggerSet == nil {
				triggerSet = &sameTypeTriggerSet{}
				triggerSets[targetType] = triggerSet
			}

			triggerSet.add(trigger)
		}
	}
	return triggerSets
}

func (s *subprobe) process(r probe.Reading) {
	if err := s.record(r); err != nil {
		// Try to stumble on. We can still send alerts.
		log.WithError(err).WithFields(log.Fields{
			"monitor":  s.monitor.id,
			"subprobe": s.name,
			"state":    r.State,
			"recorded": r.Recorded,
		}).Error("Could not record reading. Skipping saving to DB.")
	}
}

func (s *subprobe) record(r probe.Reading) error {
	return errors.Mask(s.DB.Tx(func(tx *db.Tx) error {
		stateChanged := s.state != r.State
		s.lastReading = r.Recorded
		s.state = r.State
		if stateChanged {
			s.enteredState = r.Recorded
		}
		if s.state == state.Normal {
			s.lastNormal = r.Recorded
		}
		s.saveNextReading = s.saveNextReading || stateChanged

		status := s.dbStatus()
		// TODO(eefi): Update status.Silenced.
		if err := tx.UpdateSubprobeStatus(status); err != nil {
			return errors.Maskf(err, "update subprobe status")
		}

		if s.saveNextReading {
			dbReading := db.Reading{
				SubprobeID: s.id,
				Recorded:   r.Recorded,
				State:      r.State,
			}
			if err := tx.InsertReading(dbReading); err != nil {
				return errors.Maskf(err, "insert reading")
			}

			s.saveNextReading = false
		}

		return nil
	}))
}

func (s *subprobe) dbStatus() db.SubprobeStatus {
	return db.SubprobeStatus{
		SubprobeID:   s.id,
		Recorded:     s.lastReading,
		State:        s.state,
		EnteredState: s.enteredState,
		LastNormal:   s.lastNormal,
	}
}
