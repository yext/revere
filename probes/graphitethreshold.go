package probes

import (
	"encoding/json"
	"fmt"
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
	graphite        string
	metric          string
	triggerIf       string
	thresholds      map[revere.State]float64
	holdRequirement uint
	alertFrequency  uint
	checkFrequency  uint
}

type graphiteThresholdSettings struct {
	Graphite        string
	Metric          string
	TriggerIf       string
	Thresholds      thresholds
	HoldRequirement uint
	AlertFrequency  uint
	CheckFrequency  uint
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
	gt.holdRequirement = builder.HoldRequirement
	if gt.holdRequirement == 0 {
		gt.holdRequirement = 1
	}
	gt.checkFrequency = builder.CheckFrequency
	if gt.checkFrequency == 0 {
		gt.checkFrequency = revere.CheckFrequency
	}
	gt.alertFrequency = builder.AlertFrequency
	if gt.alertFrequency == 0 {
		gt.alertFrequency = gt.checkFrequency
	}

	gt.thresholds = make(map[revere.State]float64)
	gt.thresholds[revere.Warning] = builder.Thresholds.Warning
	gt.thresholds[revere.Error] = builder.Thresholds.Error
	gt.thresholds[revere.Critical] = builder.Thresholds.Critical

	return gt, nil
}

func (gt *GraphiteThreshold) Check() (map[string]revere.Reading, error) {
	time := gt.checkFrequency * gt.holdRequirement
	resp, err := http.Get(
		fmt.Sprintf(
			"http://%s/render?target=%s&from=-%dmin&format=json",
			url.QueryEscape(gt.graphite),
			url.QueryEscape(gt.metric),
			time))
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
		var reading revere.Reading
		reading.State = revere.Normal
		state := "NORMAL"
		if trhld, ok := gt.thresholds[revere.Warning]; ok {
			if compare(gt.triggerIf, trhld, target.Datapoints) {
				reading.State = revere.Warning
				state = "WARNING"
			}
		}
		if trhld, ok := gt.thresholds[revere.Error]; ok {
			if compare(gt.triggerIf, trhld, target.Datapoints) {
				reading.State = revere.Error
				state = "ERROR"
			}
		}
		if trhld, ok := gt.thresholds[revere.Critical]; ok {
			if compare(gt.triggerIf, trhld, target.Datapoints) {
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

func compare(comparator string, threshold float64, datapoints [][]*float64) bool {
	allNil := true
	for _, datapoint := range datapoints {
		if datapoint[0] == nil {
			continue
		}
		allNil = false
		val := *datapoint[0]

		switch comparator {
		case greaterThan:
			if !(val > threshold) {
				return false
			}
		case greaterThanEq:
			if !(val >= threshold) {
				return false
			}
		case lessThan:
			if !(val < threshold) {
				return false
			}
		case lessThanEq:
			if !(val <= threshold) {
				return false
			}
		}
	}
	return !allNil
}
