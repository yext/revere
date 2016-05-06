package probes

import (
	"encoding/json"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/util"
)

type GraphiteThreshold struct{}

type GraphiteThresholdProbe struct {
	// TODO(fchen): fix tags on front-end js
	URL             string
	Expression      string
	Thresholds      ThresholdsModel
	AuditFunction   string
	CheckPeriod     int64
	CheckPeriodType string
	TriggerIf       string
	AuditPeriod     int64
	AuditPeriodType string
}

type ThresholdsModel struct {
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

func (_ GraphiteThreshold) Id() ProbeTypeId {
	return 0
}

func (_ GraphiteThreshold) Name() string {
	return "Graphite Threshold"
}

func (_ GraphiteThreshold) loadFromParams(probe string) (Probe, error) {
	var g GraphiteThresholdProbe
	err := json.Unmarshal([]byte(probe), &g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (_ GraphiteThreshold) loadFromDb(probe string) (Probe, error) {
	var g probe.GraphiteThresholdDBModel
	err := json.Unmarshal([]byte(probe), &g)
	if err != nil {
		return nil, err
	}

	checkPeriod, checkPeriodType := util.GetPeriodAndType(g.CheckPeriodMilli)
	auditPeriod, auditPeriodType := util.GetPeriodAndType(g.TimeToAuditMilli)

	return &GraphiteThresholdProbe{
		g.URL,
		g.Expression,
		ThresholdsModel{
			*g.Thresholds.Warning,
			*g.Thresholds.Error,
			*g.Thresholds.Critical,
		},
		g.AuditFunction,
		checkPeriod,
		checkPeriodType,
		g.TriggerIf,
		auditPeriod,
		auditPeriodType,
	}, nil
}

func (_ GraphiteThreshold) blank() (Probe, error) {
	return &GraphiteThresholdProbe{}, nil
}

func (_ GraphiteThreshold) Templates() map[string]string {
	return map[string]string{
		"edit": "graphite-edit.html",
		"view": "graphite-view.html",
	}
}

func (_ GraphiteThreshold) Scripts() map[string][]string {
	return map[string][]string{
		"edit": []string{
			"graphite-preview.js",
		},
	}
}

func (_ GraphiteThreshold) AcceptedDataSourceTypeIds() []datasources.DataSourceTypeId {
	return []datasources.DataSourceTypeId{
		datasources.Graphite{}.Id(),
	}
}

func (g *GraphiteThresholdProbe) Serialize() (string, error) {
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

// TODO(fchen): fix references to ProbeType() in frontend
func (g *GraphiteThresholdProbe) Type() ProbeType {
	return GraphiteThreshold{}
}

func (g *GraphiteThresholdProbe) Validate() (errs []string) {
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
