package vm

import (
	"github.com/yext/revere"
	"github.com/yext/revere/targets"
)

type Trigger struct {
	*revere.Trigger
	Target   *Target
	Subprobe string
}

type modelTrigger struct {
	Trigger  *revere.Trigger
	Subprobe string
}

func (t *Trigger) Id() int64 {
	return int64(t.Trigger.TriggerId)
}

func NewTriggersFromLabelTriggers(lts []*revere.LabelTrigger) ([]*Trigger, error) {
	triggers := make([]*modelTrigger, len(lts))
	for i, labelTrigger := range lts {
		triggers[i] = &modelTrigger{&labelTrigger.Trigger, ""}
	}

	return newTriggers(triggers)
}

func NewTriggersFromMonitorTriggers(mts []*revere.MonitorTrigger) ([]*Trigger, error) {
	triggers := make([]*modelTrigger, len(mts))
	for i, monitorTrigger := range mts {
		triggers[i] = &modelTrigger{&monitorTrigger.Trigger, monitorTrigger.Subprobe}
	}

	return newTriggers(triggers)
}

func newTrigger(t *revere.Trigger, s string) (*Trigger, error) {
	viewmodel := new(Trigger)

	viewmodel.Trigger = t
	viewmodel.Subprobe = s

	targetType, err := targets.TargetTypeById(t.TargetType)
	if err != nil {
		return nil, err
	}

	target, err := targetType.Load(t.TargetJson)
	if err != nil {
		return nil, err
	}
	viewmodel.Target = NewTarget(target)

	return viewmodel, nil
}

func newTriggers(mts []*modelTrigger) ([]*Trigger, error) {
	triggers := make([]*Trigger, len(mts))
	for i, modelTrigger := range mts {
		trigger, err := newTrigger(modelTrigger.Trigger, modelTrigger.Subprobe)
		if err != nil {
			return nil, err
		}
		triggers[i] = trigger
	}

	return triggers, nil
}

func BlankTrigger() *Trigger {
	viewmodel := new(Trigger)

	viewmodel.Trigger = new(revere.Trigger)
	viewmodel.Target = DefaultTarget()

	return viewmodel
}

func (t *Trigger) GetTargetType() targets.TargetType {
	return t.Target.TargetType()
}
