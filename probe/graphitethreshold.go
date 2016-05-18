package probe

import (
	"fmt"
	"math"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/resource"
	"github.com/yext/revere/state"
)

// GraphiteThreshold implements a probe that assigns states based on whether a
// Graphite metric is above or below various constant values.
type GraphiteThreshold struct {
	*Polling

	graphiteBase string
	expression   string
	timeToAudit  time.Duration

	thresholds      []graphiteThresholdThreshold
	summarizeValues func(values []float64) float64
	triggersOn      func(summaryValue, threshold float64) bool

	auditFunctionName string
	triggerIfText     string
}

type graphiteThresholdThreshold struct {
	state     state.State
	threshold float64
}

func newGraphiteThreshold(configJSON types.JSONText, readingsSink chan<- []Reading) (Probe, error) {
	gt := GraphiteThreshold{}

	var config GraphiteThresholdDBModel
	err := configJSON.Unmarshal(&config)
	if err != nil {
		return nil, errors.Maskf(err, "deserialize probe config")
	}

	checkPeriod := time.Duration(config.CheckPeriodMilli) * time.Millisecond
	gt.Polling, err = NewPolling(checkPeriod, &gt, readingsSink)
	if err != nil {
		return nil, errors.Mask(err)
	}

	gt.graphiteBase = fmt.Sprintf("http://%s/", config.URL)
	gt.expression = config.Expression
	gt.timeToAudit = time.Duration(config.TimeToAuditMilli) * time.Millisecond

	// Must be in increasing severity order.
	if config.Thresholds.Warning != nil {
		gt.thresholds = append(gt.thresholds,
			graphiteThresholdThreshold{state.Warning, *config.Thresholds.Warning})
	}
	if config.Thresholds.Error != nil {
		gt.thresholds = append(gt.thresholds,
			graphiteThresholdThreshold{state.Error, *config.Thresholds.Error})
	}
	if config.Thresholds.Critical != nil {
		gt.thresholds = append(gt.thresholds,
			graphiteThresholdThreshold{state.Critical, *config.Thresholds.Critical})
	}

	var ok bool

	gt.summarizeValues, ok = auditFunctions[config.AuditFunction]
	if !ok {
		return nil, errors.Errorf("unknown audit function: %s", config.AuditFunction)
	}

	gt.triggersOn, ok = triggerIfFunctions[config.TriggerIf]
	if !ok {
		return nil, errors.Errorf("unknown trigger if: %s", config.TriggerIf)
	}

	gt.auditFunctionName = config.AuditFunction
	gt.triggerIfText = config.TriggerIf

	return &gt, nil
}

var (
	auditFunctions = map[string]func([]float64) float64{
		"avg": func(values []float64) float64 {
			sum := float64(0)
			count := 0
			for _, value := range values {
				if !math.IsNaN(value) {
					sum += value
					count += 1
				}
			}
			return sum / float64(count)
		},
		"max": func(values []float64) float64 {
			max := math.Inf(-1)
			for _, value := range values {
				if value > max {
					max = value
				}
			}
			return max
		},
		"min": func(values []float64) float64 {
			min := math.Inf(+1)
			for _, value := range values {
				if value < min {
					min = value
				}
			}
			return min
		},
	}

	triggerIfFunctions = map[string]func(float64, float64) bool{
		"<": func(summaryValue, threshold float64) bool {
			return summaryValue < threshold
		},
		"<=": func(summaryValue, threshold float64) bool {
			return summaryValue <= threshold
		},
		">=": func(summaryValue, threshold float64) bool {
			return summaryValue >= threshold
		},
		">": func(summaryValue, threshold float64) bool {
			return summaryValue > threshold
		},
	}
)

func (gt *GraphiteThreshold) Check() []Reading {
	g := resource.Graphite{gt.graphiteBase}

	series, err := g.QueryRecent(gt.expression, gt.timeToAudit)
	now := time.Now()
	if err != nil {
		// TODO(eefi): Include this probe's monitor's ID.
		log.WithError(err).Error("Could not query Graphite.")

		// TODO(eefi): Put err in the details.
		return []Reading{{"_", state.Unknown, now, nil}}
	}

	if len(series) == 0 {
		return []Reading{{"_", state.Normal, now, nil}}
	}

	readings := make([]Reading, len(series)+1)
	for i, s := range series {
		r := Reading{s.Name, state.Normal, s.End, nil}

		summaryValue := gt.summarizeValues(s.Values)
		triggeredThreshold := math.NaN()
		for _, t := range gt.thresholds {
			if gt.triggersOn(summaryValue, t.threshold) {
				r.State = t.state
				triggeredThreshold = t.threshold
			}
		}

		r.Details = graphiteThresholdDetails{
			auditFunction: gt.auditFunctionName,
			timeToAudit:   gt.timeToAudit,
			triggerIf:     gt.triggerIfText,

			measured:  summaryValue,
			threshold: triggeredThreshold,

			graphite:    &g,
			expression:  gt.expression,
			seriesName:  s.Name,
			seriesStart: s.Start,
			seriesEnd:   s.End,
		}

		readings[i] = r
	}
	readings[len(series)] = Reading{"_", state.Normal, now, nil}

	return readings
}
