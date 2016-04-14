package revere_test

import (
	"database/sql"
	"testing"

	. "github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/test"
)

var (
	probeType = probes.GraphiteThreshold{}
	probeJson = test.DefaultProbeJson

	subprobes = []string{
		"[a-z]",
		"test.*\\.example",
	}
)

func validMonitor() *Monitor {
	m := new(Monitor)
	m.Name = "Name"
	m.Owner = "Test <test@example.com>"
	m.Description = "Desc."
	m.Response = "Response"
	m.ProbeType = probeType.Id()
	m.ProbeJson = probeJson
	m.Triggers = []*MonitorTrigger{validMonitorTrigger()}
	return m
}

func validMonitorTrigger() *MonitorTrigger {
	mt := new(MonitorTrigger)
	mt.Trigger = *validTrigger()
	mt.Subprobe = "test.*examples"
	mt.Delete = false
	return mt
}

func TestInvalidMonitorName(t *testing.T) {
	monitor := validMonitor()
	testDB := new(sql.DB)
	// XXX only works because of how we use the db currently,
	// should be replaced with a real dummy db eventually
	monitor.Name = ""
	if errs := monitor.Validate(testDB); errs == nil {
		t.Error("Expected error for invalid monitor name")
	}
}

func TestInvalidMonitorTriggerSubprobe(t *testing.T) {
	trigger := validMonitorTrigger()
	trigger.Subprobe = "a["
	if errs := trigger.Validate(); errs == nil {
		t.Error("Expected error for invalid subprobe")
	}
}

func TestValidTriggerSubprobe(t *testing.T) {
	trigger := validMonitorTrigger()
	for _, s := range subprobes {
		trigger.Subprobe = s
		if errs := trigger.Validate(); errs != nil {
			t.Errorf("Unexpected error for subprobe: %s\n", s)
		}
	}
}
