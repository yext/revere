package probes_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/yext/revere/probes"
	"github.com/yext/revere/test"
)

var (
	gtId        = 0
	gtName      = "Graphite Threshold"
	gtProbeType = GraphiteThreshold{}
	validGtJson = test.DefaultProbeJson
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func validGraphiteThresholdProbe() (*GraphiteThresholdProbe, error) {
	probe, err := gtProbeType.Load(validGtJson)
	if err != nil {
		return nil, err
	}

	gtProbe, ok := probe.(GraphiteThresholdProbe)
	if !ok {
		return nil, fmt.Errorf("Invalid probe loaded for probe type: %s\n", gtProbeType.Name())
	}

	return &gtProbe, nil
}

func TestGraphiteThresholdId(t *testing.T) {
	if int(gtProbeType.Id()) != gtId {
		t.Errorf("Expected graphite threshold probe type id: %d, got %d\n", gtId, gtProbeType.Id)
	}
}

func TestGraphiteThresholdName(t *testing.T) {
	if gtProbeType.Name() != gtName {
		t.Errorf("Expected graphite threshold probe type name: %s, got %s\n", gtName, gtProbeType.Name)
	}
}

func TestGraphiteThresholdEmptyLoad(t *testing.T) {
	probe, err := gtProbeType.Load(`{}`)
	if err != nil {
		t.Fatalf("Failed to load empty graphite threshold probe: %s\n", err.Error())
	}

	_, ok := probe.(GraphiteThresholdProbe)
	if !ok {
		t.Fatalf("Invalid probe loaded for probe type: %s\n", gtProbeType.Name())
	}
}

func TestInvalidGraphiteThresholdExpression(t *testing.T) {
	gtProbe, err := validGraphiteThresholdProbe()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gtProbe.Expression = ""
	errs := gtProbe.Validate()
	if errs == nil {
		t.Error("Expected error for nil graphite expression")
	}
}

func TestInvalidGraphiteThresholdPeriod(t *testing.T) {
	gtProbe, err := validGraphiteThresholdProbe()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gtProbe.CheckPeriod = -1
	errs := gtProbe.Validate()
	if errs == nil {
		t.Errorf("Expected error for invalid check period: %d\n", gtProbe.CheckPeriod)
	}
	gtProbe.CheckPeriod = 1

	gtProbe.AuditPeriod = -1
	errs = gtProbe.Validate()
	if errs == nil {
		t.Errorf("Expected error for invalid audit period: %d\n", gtProbe.AuditPeriod)
	}
	gtProbe.AuditPeriod = 1
}

func TestValidGraphiteThresholdPeriodType(t *testing.T) {
	gtProbe, err := validGraphiteThresholdProbe()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, pt := range test.PeriodTypes {
		gtProbe.CheckPeriodType = pt
		errs := gtProbe.Validate()
		if errs != nil {
			t.Errorf("Unexpected error for check period type: %s\n", pt)
		}

		gtProbe.AuditPeriodType = pt
		errs = gtProbe.Validate()
		if errs != nil {
			t.Error("Unexpected error for audit period type: %s\n", pt)
		}
	}
}

func TestInvalidGraphiteThresholdPeriodType(t *testing.T) {
	gtProbe, err := validGraphiteThresholdProbe()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gtProbe.CheckPeriodType = ""
	errs := gtProbe.Validate()
	if errs == nil {
		t.Error("Expected error for invalid check period type")
	}
	gtProbe.CheckPeriodType = test.PeriodTypes[0]

	gtProbe.AuditPeriodType = ""
	errs = gtProbe.Validate()
	if errs == nil {
		t.Error("Expected error for invalid audit period type")
	}
}
