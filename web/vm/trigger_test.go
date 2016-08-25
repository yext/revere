package vm

import (
	"testing"

	"github.com/yext/revere/state"
	"github.com/yext/revere/target"
	"github.com/yext/revere/test"
)

var (
	targetType = target.EmailType{}
	targetJson = test.DefaultTargetJson
)

func validTrigger() *Trigger {
	t := new(Trigger)
	t.Level = state.Normal
	t.Period = 5
	t.PeriodType = "minute"
	t.TriggerOnExit = false
	t.TargetType = targetType.Id()
	t.TargetParams = targetJson
	return t
}

func TestValidTriggerLevel(t *testing.T) {
	trigger := validTrigger()
	for _, s := range []state.State{state.Normal, state.Warning, state.Critical, state.Error, state.Unknown} {
		trigger.Level = s
		errs := trigger.validate()
		if errs != nil {
			t.Errorf("Unexpected error for level: %s\n", s)
		}
	}
}

func TestInvalidTriggerLevel(t *testing.T) {
	trigger := validTrigger()
	trigger.Level = 99
	errs := trigger.validate()
	if errs == nil {
		t.Error("Expected error for invalid level")
	}
}

func TestInvalidTriggerPeriod(t *testing.T) {
	trigger := validTrigger()
	trigger.PeriodType = ""
	errs := trigger.validate()
	if errs == nil {
		t.Error("Expected error for invalid period type")
	}
}

func TestValidTriggerPeriod(t *testing.T) {
	trigger := validTrigger()
	for _, s := range test.PeriodTypes {
		trigger.PeriodType = s
		errs := trigger.validate()
		if errs != nil {
			t.Errorf("Unexpected error for period type: %s\n", s)
		}
	}
}
