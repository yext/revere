package daemon

import (
	"time"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
	"github.com/yext/revere/target"
)

type trigger struct {
	level         state.State
	triggerOnExit bool
	period        time.Duration
	target        target.Target
}

func newTrigger(dbModel *db.Trigger) (*trigger, error) {
	target, err := target.New(dbModel.TargetType, dbModel.Target)
	if err != nil {
		return nil, errors.Maskf(err, "make target")
	}

	return &trigger{
		level:         dbModel.Level,
		triggerOnExit: dbModel.TriggerOnExit,
		period:        time.Duration(dbModel.PeriodMilli) * time.Millisecond,
		target:        target,
	}, nil
}

type sameTypeTriggerSet struct {
	triggers []*trigger
}

func (s *sameTypeTriggerSet) add(t *trigger) {
	s.triggers = append(s.triggers, t)
}
