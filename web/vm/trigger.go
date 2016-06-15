package vm

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
	"github.com/yext/revere/targets"
	"github.com/yext/revere/util"
)

type Trigger struct {
	TriggerID     db.TriggerID
	Level         string
	Period        int64
	PeriodType    string
	TargetType    db.TargetType
	TargetParams  string
	TriggerOnExit bool
	Target        *targets.Target
}

func newTriggerFromModel(trigger *db.Trigger) (*Trigger, error) {
	target, err := targets.LoadFromDb(trigger.TargetType, trigger.Target)
	if err != nil {
		return nil, errors.Trace(err)
	}

	period, periodType := util.GetPeriodAndType(trigger.PeriodMilli)

	return &Trigger{
		TriggerID:     trigger.TriggerID,
		Level:         trigger.Level,
		Period:        period,
		PeriodType:    periodType,
		TargetType:    trigger.TargetType,
		TargetParams:  nil,
		TriggerOnExit: trigger.TriggerOnExit,
		Target:        target,
	}, nil
}

func (t *Trigger) Id() int64 {
	return int64(t.TriggerID)
}

func (t *Trigger) validate() (errs []string) {
	var err error
	t.Target, err = targets.LoadFromParams(t.TargetType, t.TargetParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load probe for monitor: %s", m.ProbeParams))
	}
	errs = append(errs, target.Validate()...)

	errs = append(errs, m.Probe.Validate()...)
	if _, err := state.ReverseStates(t.Level); err != nil {
		errs = append(errs, fmt.Sprintf("Invalid state for trigger: %s", t.Level))
	}

	if util.GetMs(t.Period, t.PeriodType) == 0 {
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

	return &db.Trigger{
		TriggerID:     t.TriggerID,
		Level:         t.Level,
		TriggerOnExit: t.TriggerOnExit,
		PeriodMilli:   util.GetMs(t.Period, t.PeriodType),
		TargetType:    t.TargetType,
		Target:        triggerJson,
	}, nil
}
