package daemon

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
	"github.com/yext/revere/target"
)

type triggerTemplate struct {
	id            db.TriggerID
	level         state.State
	triggerOnExit bool
	period        time.Duration
	target        target.Target
}

func newTriggerTemplate(dbModel *db.Trigger) (*triggerTemplate, error) {
	target, err := target.New(dbModel.TargetType, dbModel.Target)
	if err != nil {
		return nil, errors.Maskf(err, "make target")
	}

	return &triggerTemplate{
		level:         dbModel.Level,
		triggerOnExit: dbModel.TriggerOnExit,
		period:        time.Duration(dbModel.PeriodMilli) * time.Millisecond,
		target:        target,
	}, nil
}

type trigger struct {
	*triggerTemplate
	lastAlert time.Time
}

func newTrigger(template *triggerTemplate) *trigger {
	return &trigger{triggerTemplate: template}
}

func (t *trigger) shouldTrigger(a *target.Alert) bool {
	// TODO(eefi): Implement.
	return true
}

type sameTypeTriggerSet map[db.TriggerID]*trigger

func newSameTypeTriggerSet() sameTypeTriggerSet {
	return make(map[db.TriggerID]*trigger)
}

func (s sameTypeTriggerSet) add(t *trigger) {
	s[t.id] = t
}

func (s sameTypeTriggerSet) alert(a *target.Alert) {
	toAlert := make(map[db.TriggerID]target.Target)
	var inactive []target.Target
	var targetType target.Type
	for _, trigger := range s {
		if trigger.shouldTrigger(a) {
			toAlert[trigger.id] = trigger.target
			targetType = trigger.target.Type()
		} else {
			inactive = append(inactive, trigger.target)
		}
	}

	if len(toAlert) == 0 {
		return
	}

	errors := targetType.Alert(a, toAlert, inactive)

	for _, errAndIDs := range errors {
		log.WithError(errAndIDs.Err).WithFields(log.Fields{
			"monitor":    a.MonitorID,
			"subprobe":   a.SubprobeName,
			"state":      a.NewState,
			"recorded":   a.Recorded,
			"targetType": targetType.ID(),
			"triggers":   errAndIDs.IDs,
		}).Error("Some alerts failed and have been lost.")

		for _, id := range errAndIDs.IDs {
			delete(toAlert, id)
		}
	}

	now := time.Now()
	for id, _ := range toAlert {
		s[id].lastAlert = now
	}
}
