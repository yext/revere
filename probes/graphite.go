package probes

import (
	"encoding/json"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/util"
)

type GraphiteThreshold struct{}

type GraphiteThresholdProbe struct {
	Url        string `json:"url"`
	Expression string `json:"expression"`
	Threshold
	AuditFunction   string   `json:"auditFunction"`
	CheckPeriod     int64    `json:"checkPeriod"`
	CheckPeriodType string   `json:"checkPeriodType"`
	TriggerIf       string   `json:"triggerIf"`
	AuditPeriod     int64    `json:"auditPeriod"`
	AuditPeriodType string   `json:"auditPeriodType"`
	DataSources     []string `json:"-"`
}

type Threshold struct {
	Warning  int64 `json:"warningThreshold"`
	Error    int64 `json:"errorThreshold"`
	Critical int64 `json:"criticalThreshold"`
}

var (
	templates = map[string]string{
		"edit": "graphite-edit.html",
		"view": "graphite-view.html",
	}
	scripts = map[string][]string{
		"edit": []string{
			"graphite-preview.js",
		},
	}
	AcceptedDataSourceTypes = []datasources.DataSourceTypeId{
		datasources.GraphiteDataSource{}.Id(),
	}

	// All graphite datasources found in the conf file
	GraphiteUrls []string
)

func init() {
	addProbeType(GraphiteThreshold{})
}

func SetGraphiteUrls(urls []string) {
	GraphiteUrls = urls
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
		return g, err
	}
	g.DataSources = GraphiteUrls
	return g, nil
}

func (gt GraphiteThreshold) Templates() map[string]string {
	return templates
}

func (gt GraphiteThreshold) Scripts() map[string][]string {
	return scripts
}

func (g GraphiteThreshold) AcceptedDataSourceTypeIds() []datasources.DataSourceTypeId {
	return AcceptedDataSourceTypes
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
	if util.GetMs(g.CheckPeriod, g.CheckPeriodType) == 0 {
		errs = append(errs, "Invalid check period")
	}

	if util.GetMs(g.AuditPeriod, g.AuditPeriodType) == 0 {
		errs = append(errs, "Invalid alert period")
	}

	return
}
