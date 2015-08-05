package probes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/yext/revere"
)

const (
	greaterThan   string = ">"
	greaterThanEq string = ">="
	lessThan      string = "<"
	lessThanEq    string = "<="
)

// A GraphiteThreshold probe compares Graphite metrics to static values.
type GraphiteThreshold struct {
	graphite   string
	metric     string
	triggerIf  string
	thresholds map[revere.State]float64
}

type graphiteThresholdSettings struct {
	Graphite   string
	Metric     string
	TriggerIf  string
	Thresholds thresholds
}

type thresholds struct {
	Warning  float64
	Error    float64
	Critical float64
}

type graphiteData struct {
	Target     string
	Datapoints [][]*float64
}

type graphiteReadingDetails struct {
	details string
}

func (r graphiteReadingDetails) Text() string {
	return r.details
}

// NewGraphiteThreshold builds a probe based on serialized settings.
//
// TODO(eefi): Detail the serialization format.
func NewGraphiteThreshold(settings string) (*GraphiteThreshold, error) {
	builder := new(graphiteThresholdSettings)
	err := json.Unmarshal([]byte(settings), &builder)
	if err != nil {
		return nil, err
	}

	gt := new(GraphiteThreshold)
	gt.graphite = builder.Graphite
	gt.metric = builder.Metric
	gt.triggerIf = builder.TriggerIf

	gt.thresholds = make(map[revere.State]float64)
	gt.thresholds[revere.Warning] = builder.Thresholds.Warning
	gt.thresholds[revere.Error] = builder.Thresholds.Error
	gt.thresholds[revere.Critical] = builder.Thresholds.Critical

	return gt, nil
}

func (gt *GraphiteThreshold) Check() (map[string]revere.Reading, error) {
	resp, err := http.Get(
		"http://" +
			url.QueryEscape(gt.graphite) +
			"/render?target=" +
			url.QueryEscape(gt.metric) +
			"&from=-5min" +
			"&format=json")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []graphiteData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var readings = make(map[string]revere.Reading)
	for _, target := range data {
		val := 0.0
		for i := len(target.Datapoints) - 1; i >= 0; i-- {
			if target.Datapoints[i][0] != nil {
				val = *target.Datapoints[i][0]
			}
		}

		var reading revere.Reading
		reading.State = revere.Normal
		state := "NORMAL"
		if trhld, ok := gt.thresholds[revere.Warning]; ok {
			if compare(gt.triggerIf, trhld, val) {
				reading.State = revere.Warning
				state = "WARNING"
			}
		}
		if trhld, ok := gt.thresholds[revere.Error]; ok {
			if compare(gt.triggerIf, trhld, val) {
				reading.State = revere.Error
				state = "ERROR"
			}
		}
		if trhld, ok := gt.thresholds[revere.Critical]; ok {
			if compare(gt.triggerIf, trhld, val) {
				reading.State = revere.Critical
				state = "CRITICAL"
			}
		}

		reading.Details = graphiteReadingDetails{
			"Target " +
				target.Target +
				" has state " +
				state +
				" for metric " +
				gt.metric}
		readings[target.Target] = reading
	}
	return readings, nil
}

func compare(comparator string, threshold float64, val float64) bool {
	switch comparator {
	case greaterThan:
		if val > threshold {
			return true
		}
	case greaterThanEq:
		if val >= threshold {
			return true
		}
	case lessThan:
		if val < threshold {
			return true
		}
	case lessThanEq:
		if val <= threshold {
			return true
		}
	}
	return false
}
