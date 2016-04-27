package probes

import (
	"encoding/json"

	"github.com/yext/revere/datasources"
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

type Threshold struct {
	Warning  int64
	Error    int64
	Critical int64
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
	return 0
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
