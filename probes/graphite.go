package probes

import (
	"encoding/json"
	"strconv"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/util"
)

type GraphiteThreshold struct{}

type GraphiteThresholdProbe struct {
	GraphiteThreshold

	// TODO(fchen): fix tags on front-end js
	URL               string
	SourceID          db.DatasourceID
	Expression        string
	Thresholds        ThresholdsModel
	AuditFunction     string
	CheckPeriod       int64
	CheckPeriodType   string
	TriggerIf         string
	AuditPeriod       int64
	AuditPeriodType   string
	IgnoredPeriod     int64
	IgnoredPeriodType string
}

type ThresholdsModel struct {
	Warning  *float64
	Error    *float64
	Critical *float64
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

func (GraphiteThreshold) Id() db.ProbeType {
	return 1
}

func (GraphiteThreshold) Name() string {
	return "Graphite Threshold"
}

func (GraphiteThreshold) loadFromParams(probe string) (Probe, error) {
	var g GraphiteThresholdProbe
	err := json.Unmarshal([]byte(probe), &g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (GraphiteThreshold) loadFromDb(encodedProbe string, tx *db.Tx) (Probe, error) {
	var g probe.GraphiteThresholdDBModel
	err := json.Unmarshal([]byte(encodedProbe), &g)
	if err != nil {
		return nil, err
	}

	checkPeriod, checkPeriodType := util.GetPeriodAndType(g.CheckPeriodMilli)
	auditPeriod, auditPeriodType := util.GetPeriodAndType(g.TimeToAuditMilli)
	ignoredPeriod, ignoredPeriodType := util.GetPeriodAndType(g.RecentTimeToIgnoreMilli)

	// Ignored for now - will swap next commit, when the datasource is updated
	_, err = tx.LoadDatasource(db.DatasourceID(g.SourceID))

	return &GraphiteThresholdProbe{
		URL:        g.URL,
		SourceID:   db.DatasourceID(g.SourceID),
		Expression: g.Expression,
		Thresholds: ThresholdsModel{
			g.Thresholds.Warning,
			g.Thresholds.Error,
			g.Thresholds.Critical,
		},
		AuditFunction:     g.AuditFunction,
		CheckPeriod:       checkPeriod,
		CheckPeriodType:   checkPeriodType,
		TriggerIf:         g.TriggerIf,
		AuditPeriod:       auditPeriod,
		AuditPeriodType:   auditPeriodType,
		IgnoredPeriod:     ignoredPeriod,
		IgnoredPeriodType: ignoredPeriodType,
	}, nil
}

func (GraphiteThreshold) blank() (Probe, error) {
	return &GraphiteThresholdProbe{}, nil
}

func (GraphiteThreshold) Templates() map[string]string {
	return map[string]string{
		"edit": "graphite-edit.html",
		"view": "graphite-view.html",
	}
}

func (gt GraphiteThreshold) Scripts() map[string][]string {
	return map[string][]string{
		"edit": []string{
			"graphite-threshold.js",
			"graphite-ds-loader.js",
			gt.PreviewScript(),
		},
	}
}

func (GraphiteThreshold) PreviewScript() string {
	return "graphite-preview.js"
}

func (GraphiteThreshold) AcceptedSourceTypes() []db.SourceType {
	return []db.SourceType{
		datasources.Graphite{}.Id(),
	}
}

func (g GraphiteThresholdProbe) PreviewParams() map[string]string {
	var warningStr, errorStr, criticalStr string
	if g.Thresholds.Warning != nil {
		warningStr = strconv.FormatFloat(*g.Thresholds.Warning, 'f', -1, 64)
	}
	if g.Thresholds.Error != nil {
		errorStr = strconv.FormatFloat(*g.Thresholds.Error, 'f', -1, 64)
	}
	if g.Thresholds.Critical != nil {
		criticalStr = strconv.FormatFloat(*g.Thresholds.Critical, 'f', -1, 64)
	}
	return map[string]string{
		"Expression": g.Expression,
		"URL":        g.URL,
		"Warning":    warningStr,
		"Error":      errorStr,
		"Critical":   criticalStr,
	}
}

func (g GraphiteThresholdProbe) Serialize() (string, error) {
	checkPeriodMilli := util.GetMs(g.CheckPeriod, g.CheckPeriodType)
	auditPeriodMilli := util.GetMs(g.AuditPeriod, g.AuditPeriodType)
	ignoredPeriodMilli := util.GetMs(g.IgnoredPeriod, g.IgnoredPeriodType)

	gtDB := probe.GraphiteThresholdDBModel{
		URL:        g.URL,
		SourceID:   int64(g.SourceID),
		Expression: g.Expression,
		Thresholds: probe.GraphiteThresholdThresholdsDBModel{
			Warning:  g.Thresholds.Warning,
			Error:    g.Thresholds.Error,
			Critical: g.Thresholds.Critical,
		},
		TriggerIf:               g.TriggerIf,
		CheckPeriodMilli:        checkPeriodMilli,
		TimeToAuditMilli:        auditPeriodMilli,
		RecentTimeToIgnoreMilli: ignoredPeriodMilli,
		AuditFunction:           g.AuditFunction,
	}

	gtDBJSON, err := json.Marshal(gtDB)
	return string(gtDBJSON), err
}

// TODO(fchen): fix references to ProbeType() in frontend
func (g GraphiteThresholdProbe) Type() ProbeType {
	return GraphiteThreshold{}
}

func (g GraphiteThresholdProbe) Validate() (errs []string) {
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
