package revere_test

import (
	"testing"

	. "github.com/yext/revere"
	"github.com/yext/revere/targets"
	"github.com/yext/revere/test"
)

var (
	targetType = targets.Email{}
	targetJson = test.DefaultTargetJson

	subprobes = []string{
		"[a-z]",
		"test.*.example",
	}
)

func validTrigger() *Trigger {
	t := new(Trigger)
	t.Level = States(NORMAL)
	t.Period = 5
	t.PeriodType = "minute"
	t.Subprobe = "test.*.examples"
	t.TriggerOnExit = false
	t.TargetType = targetType.Id()
	t.TargetJson = targetJson
	return t
}

func TestValidTriggerLevel(t *testing.T) {
	trigger := validTrigger()
	for s, _ := range ReverseStates {
		trigger.Level = s
		errs := trigger.Validate()
		if errs != nil {
			t.Errorf("Unexpected error for level: %s\n", s)
		}
	}
}

func TestInvalidTriggerLevel(t *testing.T) {
	trigger := validTrigger()
	trigger.Level = ""
	errs := trigger.Validate()
	if errs == nil {
		t.Error("Expected error for invalid level")
	}
}

func TestInvalidTriggerPeriod(t *testing.T) {
	trigger := validTrigger()
	trigger.PeriodType = ""
	errs := trigger.Validate()
	if errs == nil {
		t.Error("Expected error for invalid period type")
	}
}

func TestValidTriggerPeriod(t *testing.T) {
	trigger := validTrigger()
	for _, s := range test.PeriodTypes {
		trigger.PeriodType = s
		errs := trigger.Validate()
		if errs != nil {
			t.Errorf("Unexpected error for period type: %s\n", s)
		}
	}
}

func TestInvalidTriggerSubprobe(t *testing.T) {
	trigger := validTrigger()
	trigger.Subprobe = ""
	errs := trigger.Validate()
	if errs == nil {
		t.Error("Expected error for invalid subprobe")
	}
}

func TestValidTriggerSubprobe(t *testing.T) {
	trigger := validTrigger()
	for _, s := range subprobes {
		trigger.Subprobe = s
		errs := trigger.Validate()
		if errs != nil {
			t.Errorf("Unexpected error for subprobe: %s\n", s)
		}
	}
}
