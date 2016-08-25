package vm

import (
	"testing"

	"github.com/yext/revere/db"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/test"
)

var (
	probeType = probe.GraphiteThresholdType{}
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
	m.ProbeParams = probeJson
	m.Triggers = []*MonitorTrigger{validMonitorTrigger()}
	return m
}

func validMonitorTrigger() *MonitorTrigger {
	mt := &MonitorTrigger{}
	mt.Subprobes = "test.*examples"
	mt.Trigger = BlankTrigger()
	return mt
}

func TestInvalidMonitorName(t *testing.T) {
	monitor := validMonitor()
	testDB := new(db.DB)
	// XXX only works because of how we use the db currently,
	// should be replaced with a real dummy db eventually
	monitor.Name = ""
	if errs := monitor.Validate(testDB); errs == nil {
		t.Error("Expected error for invalid monitor name")
	}
}

func TestInvalidMonitorTriggerSubprobe(t *testing.T) {
	if errs := validateSubprobeRegex("a["); errs == nil {
		t.Error("Expected error for invalid subprobe")
	}
}

func TestValidTriggerSubprobe(t *testing.T) {
	for _, s := range subprobes {
		if errs := validateSubprobeRegex(s); errs != nil {
			t.Errorf("Unexpected error for subprobe: %s\n", s)
		}
	}
}
