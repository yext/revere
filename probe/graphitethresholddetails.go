package probe

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/yext/revere/durationfmt"
	"github.com/yext/revere/resource"
)

type graphiteThresholdDetails struct {
	auditFunction string
	timeToAudit   time.Duration
	triggerIf     string

	measured  float64
	threshold float64

	graphite    *resource.Graphite
	expression  string
	seriesName  string
	seriesStart time.Time
	seriesEnd   time.Time
}

func (d graphiteThresholdDetails) Text() string {
	timeToAuditText := durationfmt.ExactMulti().Format(d.timeToAudit)
	measuredText := fmt.Sprintf("%s of last %s", d.auditFunction, timeToAuditText)

	var thresholdText, thresholdVal string
	if math.IsNaN(d.threshold) {
		thresholdText, thresholdVal = "", ""
	} else {
		thresholdText = fmt.Sprintf(" %s threshold", d.triggerIf)
		thresholdVal = fmt.Sprintf(" %s %g", d.triggerIf, d.threshold)
	}

	firstLine := fmt.Sprintf("%s%s: %g%s",
		measuredText, thresholdText, d.measured, thresholdVal)

	return fmt.Sprintf("%s\n\nGraph: %s\nValues: %s\n", firstLine, d.graphURL(), d.valuesURL())
}

func (d graphiteThresholdDetails) graphURL() string {
	measuredStart := d.seriesEnd.Add(-d.timeToAudit)

	targets := make([]string, 0, 3)

	timeHighlight := fmt.Sprintf(
		`color(drawAsInfinite(timeSlice(timeFunction("", 1), "%s", "%s")), "yellow")`,
		resource.GraphiteTimestamp(measuredStart),
		resource.GraphiteTimestamp(d.seriesEnd))
	targets = append(targets, timeHighlight)

	if !math.IsNaN(d.threshold) {
		thresholdLine := fmt.Sprintf(`color(constantLine(%g), "red")`, d.threshold)
		targets = append(targets, thresholdLine)
	}

	targets = append(targets, fmt.Sprintf(`color(%s, "green")`, d.target()))

	args := map[string]string{
		"hideLegend": "true",
		"width":      "970",
		"height":     "600",
		"tz":         "UTC",
	}

	title := d.seriesName
	if !math.IsNaN(d.threshold) {
		title += fmt.Sprintf(" vs %g", d.threshold)
	}
	args["title"] = title

	contextTime := d.timeToAudit
	if contextTime < 30*time.Minute {
		contextTime = 30 * time.Minute
	}
	args["from"] = resource.GraphiteTimestamp(measuredStart.Add(-2 * contextTime))
	args["until"] = resource.GraphiteTimestamp(d.seriesEnd.Add(contextTime))

	return d.graphite.RenderURL(targets, args)
}

func (d graphiteThresholdDetails) valuesURL() string {
	return d.graphite.RenderURL([]string{d.target()}, map[string]string{
		"format": "csv",
		"tz":     "UTC",
		"from":   resource.GraphiteTimestamp(d.seriesStart),
		"until":  resource.GraphiteTimestamp(d.seriesEnd),
	})
}

func (d graphiteThresholdDetails) target() string {
	return fmt.Sprintf(`grep(%s, "^%s$")`, d.expression, regexp.QuoteMeta(d.seriesName))
}
