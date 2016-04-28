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

	return &subprobe{
		id:              status.SubprobeID,
		monitor:         monitor,
		name:            name,
		lastReading:     status.Recorded,
		state:           status.State,
		enteredState:    status.EnteredState,
		lastNormal:      status.LastNormal,
		saveNextReading: false,
		triggerSets:     triggerSets,
		Env:             monitor.Env,
	}
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

		status := db.SubprobeStatus{
			SubprobeID:   s.id,
			Recorded:     s.lastReading,
			State:        s.state,
			Silenced:     false, // TODO(eefi)
			EnteredState: s.enteredState,
			LastNormal:   s.lastNormal,
		}
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
