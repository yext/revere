package probes

import (
	"encoding/json"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/util"
)

type GraphiteThreshold struct{}

type GraphiteThresholdProbe struct {
	Url        string
	Expression string
	Threshold
	AuditFunction   string
	CheckPeriod     int64
	CheckPeriodType string
	TriggerIf       string
	AuditPeriod     int64
	AuditPeriodType string
}

// TODO(psingh): Make into *float64 and handle showing <nil> in frontend
type Threshold struct {
	Warning  float64
	Error    float64
	Critical float64
}

var (
	validGraphitePeriodTypes = []string{
		"day",
		"hour",
		"minute",
		"second",
	}
)

func init() {
	addProbeType(GraphiteThreshold{})
}

func (gt GraphiteThreshold) Id() ProbeTypeId {
	return 1
}

func (gt GraphiteThreshold) Name() string {
	return "Graphite Threshold"
}

func (gt GraphiteThreshold) Load(probe string) (Probe, error) {
	var g GraphiteThresholdProbe
	err := json.Unmarshal([]byte(probe), &g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// TODO(psingh): Clean up, temporary fix
func (gt GraphiteThreshold) LoadFromDB(probeJSON string) (Probe, error) {
	if probeJSON == `{}` {
		return &GraphiteThresholdProbe{}, nil
	}

	var g probe.GraphiteThresholdDBModel
	err := json.Unmarshal([]byte(probeJSON), &g)
	if err != nil {
		return nil, err
	}

	checkPeriod, checkPeriodType := util.GetPeriodAndType(g.CheckPeriodMilli)
	auditPeriod, auditPeriodType := util.GetPeriodAndType(g.TimeToAuditMilli)

	return &GraphiteThresholdProbe{
		g.URL,
		g.Expression,
		Threshold{
			Warning:  *g.Thresholds.Warning,
			Error:    *g.Thresholds.Error,
			Critical: *g.Thresholds.Critical,
		},
		g.AuditFunction,
		checkPeriod,
		checkPeriodType,
		g.TriggerIf,
		auditPeriod,
		auditPeriodType,
	}, nil
}

func (gt GraphiteThreshold) Templates() map[string]string {
	return map[string]string{
		"edit": "graphite-edit.html",
		"view": "graphite-view.html",
	}
}

func (gt GraphiteThreshold) Scripts() map[string][]string {
	return map[string][]string{
		"edit": []string{
			"graphite-preview.js",
		},
	}
}

func (g GraphiteThreshold) AcceptedDataSourceTypeIds() []datasources.DataSourceTypeId {
	return []datasources.DataSourceTypeId{
		datasources.Graphite{}.Id(),
	}
}

func (g GraphiteThresholdProbe) DBModelJSON() (string, error) {
	checkPeriodMilli := util.GetMs(g.CheckPeriod, g.CheckPeriodType)
	auditPeriodMilli := util.GetMs(g.AuditPeriod, g.AuditPeriodType)

	gtDB := probe.GraphiteThresholdDBModel{
		g.Url,
		g.Expression,
		probe.GraphiteThresholdThresholdsDBModel{
			Warning:  &g.Warning,
			Error:    &g.Error,
			Critical: &g.Critical,
		},
		g.TriggerIf,
		checkPeriodMilli,
		auditPeriodMilli,
		g.AuditFunction,
	}

	gtDBJSON, err := json.Marshal(gtDB)
	return string(gtDBJSON), err
}

func (g GraphiteThresholdProbe) ProbeType() ProbeType {
	return GraphiteThreshold{}
}

func (g GraphiteThresholdProbe) Validate() (errs []string) {
	if g.Url == "" {
		errs = append(errs, "Graphite data source is required")
	}

	if g.Expression == "" {
		errs = append(errs, "Graphite expression is required")
	}

	isValidCheckPeriodType := false
	for _, vpt := range validGraphitePeriodTypes {
		if g.CheckPeriodType == vpt {
			isValidCheckPeriodType = true
			break
		}
	}
	if !isValidCheckPeriodType {
		errs = append(errs, "Invalid check period type")
	}

	isValidAuditPeriodType := false
	for _, vpt := range validGraphitePeriodTypes {
		if g.AuditPeriodType == vpt {
			isValidAuditPeriodType = true
			break
		}
	}
	if !isValidAuditPeriodType {
		errs = append(errs, "Invalid audit period type")
	}

	if util.GetMs(g.CheckPeriod, g.CheckPeriodType) <= 0 {
		errs = append(errs, "Invalid check period")
	}

	if util.GetMs(g.AuditPeriod, g.AuditPeriodType) <= 0 {
		errs = append(errs, "Invalid audit period")
	}

	return
}
