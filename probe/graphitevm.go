package probe

import (
	"encoding/json"
	"strconv"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/resource"
	"github.com/yext/revere/util"
)

type GraphiteThresholdType struct{}

type GraphiteThresholdProbe struct {
	GraphiteThresholdType

	URL               string
	ResourceID        db.ResourceID
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
	addType(GraphiteThresholdType{})
}

func (GraphiteThresholdType) Id() db.ProbeType {
	return 1
}

func (GraphiteThresholdType) Name() string {
	return "Graphite Threshold"
}

func (GraphiteThresholdType) loadFromParams(probe string) (VM, error) {
	var g GraphiteThresholdProbe
	err := json.Unmarshal([]byte(probe), &g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (GraphiteThresholdType) loadFromDb(encodedProbe string, tx *db.Tx) (VM, error) {
	var g GraphiteThresholdDBModel
	err := json.Unmarshal([]byte(encodedProbe), &g)
	if err != nil {
		return nil, err
	}

	checkPeriod, checkPeriodType := util.GetPeriodAndType(g.CheckPeriodMilli)
	auditPeriod, auditPeriodType := util.GetPeriodAndType(g.TimeToAuditMilli)
	ignoredPeriod, ignoredPeriodType := util.GetPeriodAndType(g.RecentTimeToIgnoreMilli)

	dbds, err := tx.LoadResource(db.ResourceID(g.ResourceID))
	if err != nil {
		return nil, err
	}

	if dbds == nil {
		return nil, errors.Errorf("no resource found: %d")
	}

	ds, err := resource.LoadFromDB(resource.GraphiteResource{}.Id(), dbds.Resource)
	if err != nil {
		return nil, err
	}

	gds, found := ds.(*resource.GraphiteResource)
	if !found {
		return nil, errors.New("not a graphite resource")
	}

	return &GraphiteThresholdProbe{
		URL:        gds.URL,
		ResourceID: db.ResourceID(g.ResourceID),
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

func (GraphiteThresholdType) blank() (VM, error) {
	return &GraphiteThresholdProbe{}, nil
}

func (GraphiteThresholdType) Templates() map[string]string {
	return map[string]string{
		"edit": "graphite-edit.html",
		"view": "graphite-view.html",
	}
}

func (gt GraphiteThresholdType) Scripts() map[string][]string {
	return map[string][]string{
		"edit": []string{
			"graphite-threshold.js",
			"graphite-resource-loader.js",
		},
		"preview": []string{
			"graphite-preview.js",
		},
	}
}

func (GraphiteThresholdType) AcceptedResourceTypes() []db.ResourceType {
	return []db.ResourceType{
		resource.Graphite{}.Id(),
	}
}

func (g GraphiteThresholdProbe) HasResource(id db.ResourceID) bool {
	return g.ResourceID == id
}

func (g GraphiteThresholdProbe) SerializeForFrontend() map[string]string {
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

func (g GraphiteThresholdProbe) SerializeForDB() (string, error) {
	checkPeriodMilli := util.GetMs(g.CheckPeriod, g.CheckPeriodType)
	auditPeriodMilli := util.GetMs(g.AuditPeriod, g.AuditPeriodType)
	ignoredPeriodMilli := util.GetMs(g.IgnoredPeriod, g.IgnoredPeriodType)

	gtDB := GraphiteThresholdDBModel{
		ResourceID: int64(g.ResourceID),
		Expression: g.Expression,
		Thresholds: GraphiteThresholdThresholdsDBModel{
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

func (g GraphiteThresholdProbe) Type() VMType {
	return GraphiteThresholdType{}
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
