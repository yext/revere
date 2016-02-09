package revere_test

import (
	"testing"

	. "github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/test"
)

var (
	probeType = probes.GraphiteThreshold{}
	probeJson = test.DefaultProbeJson
)

func validMonitor() *Monitor {
	m := new(Monitor)
	m.Name = "Name"
	m.Owner = "Test <test@example.com>"
	m.Description = "Desc."
	m.Response = "Response"
	m.ProbeType = probeType.Id()
	m.ProbeJson = probeJson
	m.Triggers = []*Trigger{validTrigger()}
	return m
}

func TestInvalidMonitorName(t *testing.T) {
	monitor := validMonitor()
	monitor.Name = ""
	errs := monitor.Validate()
	if errs == nil {
		t.Error("Expected error for invalid monitor name")
	}
}
