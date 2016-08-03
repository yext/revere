package vm

import (
	"fmt"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
	"github.com/yext/revere/target"
	"github.com/yext/revere/util"
)

type Trigger struct {
	TriggerID     db.TriggerID
	Level         state.State
	LevelText     string
	Period        int32
	PeriodType    string
	TargetType    db.TargetType
	TargetParams  string
	TriggerOnExit bool
	Target        target.VM
	Delete        bool
}

func newTriggerFromModel(trigger *db.Trigger) (*Trigger, error) {
	target, err := target.LoadFromDb(trigger.TargetType, string(trigger.Target))
	if err != nil {
		return nil, errors.Trace(err)
	}

	period, periodType := util.GetPeriodAndType(int64(trigger.PeriodMilli))

	return &Trigger{
		TriggerID:     trigger.TriggerID,
		Level:         trigger.Level,
		Period:        int32(period),
		PeriodType:    periodType,
		TargetType:    trigger.TargetType,
		TargetParams:  "",
		TriggerOnExit: trigger.TriggerOnExit,
		Target:        target,
	}, nil
}

func BlankTrigger() *Trigger {
	return &Trigger{
		Target: target.Default(),
	}
}

func (t *Trigger) Id() int64 {
	return int64(t.TriggerID)
}

func (t *Trigger) validate() (errs []string) {
	var err error
	target, err := target.LoadFromParams(t.TargetType, t.TargetParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load target for trigger: %s", t.TargetParams))
	}
	t.Target = target
	errs = append(errs, target.Validate()...)

	if err = t.Level.Validate(); err != nil {
		errs = append(errs, fmt.Sprintf("Invalid state for trigger: %d", t.Level))
	}

	if util.GetMs(int64(t.Period), t.PeriodType) == 0 {
		errs = append(errs, fmt.Sprintf("Invalid period for trigger: %d %s", t.Period, t.PeriodType))
	}

	return
}

func (t *Trigger) setId(id db.TriggerID) {
	t.TriggerID = id
}

func (t *Trigger) toDBTrigger() (*db.Trigger, error) {
	triggerJSON, err := t.Target.Serialize()
	if err != nil {
		return nil, errors.Trace(err)
	}

	level, err := state.FromString(t.LevelText)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &db.Trigger{
		TriggerID:     t.TriggerID,
		Level:         level,
		TriggerOnExit: t.TriggerOnExit,
		PeriodMilli:   int32(util.GetMs(int64(t.Period), t.PeriodType)),
		TargetType:    t.TargetType,
		Target:        types.JSONText(triggerJSON),
	}, nil
}
