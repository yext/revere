package probes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yext/revere"
)

const (
	greaterThan   string = ">"
	greaterThanEq string = ">="
	lessThan      string = "<"
	lessThanEq    string = "<="

	graphiteUrlFormat string = "http://%s/render?target=%s&from=-%dmin&format=json"
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
	var builder graphiteThresholdSettings
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
		// alertFrequency is in seconds
		gt.alertFrequency = gt.checkFrequency * 60
	}

	gt.thresholds = make(map[revere.State]float64)
	gt.thresholds[revere.Warning] = builder.Thresholds.Warning
	gt.thresholds[revere.Error] = builder.Thresholds.Error
	gt.thresholds[revere.Critical] = builder.Thresholds.Critical

	return gt, nil
}

func newGraphiteThresholdSettings(graphite, metric, triggerIf string, w, e, c float64, hold, alertFreq, checkFreq uint) *graphiteThresholdSettings {
	gt := new(graphiteThresholdSettings)
	gt.Graphite = graphite
	gt.Metric = metric
	gt.TriggerIf = triggerIf

	gt.Thresholds = thresholds{w, e, c}

	gt.HoldRequirement = hold
	if gt.HoldRequirement == 0 {
		gt.HoldRequirement = 1
	}
	gt.CheckFrequency = checkFreq
	if gt.CheckFrequency == 0 {
		gt.CheckFrequency = revere.CheckFrequency
	}
	gt.AlertFrequency = alertFreq
	if gt.AlertFrequency == 0 {
		// alertFrequency is in seconds
		gt.AlertFrequency = gt.CheckFrequency * 60
	}

	return gt
}

// Validate accepts the url values from the http form then returns the string that
// should be stored as the JSON configuration
func Validate(form url.Values) (conf string, errs []error) {
	graphite := form.Get("graphite")
	if len(graphite) == 0 {
		errs = append(errs, errors.New("Graphite URL should not be empty"))
	}
	metric := form.Get("metric")
	if len(metric) == 0 {
		errs = append(errs, errors.New("Target Metric should not be empty"))
	}
	if len(graphite) > 0 && len(metric) > 0 {
		graphiteUrl := fmt.Sprintf(graphiteUrlFormat, graphite, metric, 1)
		resp, err := http.Get(graphiteUrl)
		if err != nil {
			errs = append(errs, errors.New("Unable to connect to graphite at this time: "+err.Error()))
		}
		if resp != nil && resp.StatusCode != 200 {
			errs = append(errs, errors.New("Graphite target appears invalid with Url: "+graphiteUrl+" status: "+resp.Status))
		}
	}

	triggerIf := form.Get("triggerIf")
	triggers := map[string]bool{
		">":  true,
		">=": true,
		"<=": true,
		"<":  true,
	}
	if !triggers[triggerIf] {
		t := make([]string, len(triggers))
		i := 0
		for k := range triggers {
			t[i] = k
			i++
		}
		errs = append(errs, errors.New("Trigger: "+triggerIf+" should be one of "+strings.Join(t, ", ")))
	}

	i, err, _ := checkInt("Hold requirement", form.Get("holdRequirement"))
	var hold uint
	if err != nil {
		errs = append(errs, err)
	} else {
		hold = i
	}
	i, err, _ = checkInt("Alert frequency", form.Get("alertFrequency"))
	var alertFreq uint
	if err != nil {
		errs = append(errs, err)
	} else {
		alertFreq = i
	}
	i, err, _ = checkInt("Check frequency", form.Get("checkFrequency"))
	var checkFreq uint
	if err != nil {
		errs = append(errs, err)
	} else {
		checkFreq = i
	}

	f, err, ok := checkFloat("Warning threshold", form.Get("threshold.warning"))
	var w float64
	if err != nil {
		errs = append(errs, err)
	} else if ok {
		w = f
	} else {
		errs = append(errs, errors.New("Warning threshold should not be empty"))
	}
	f, err, ok = checkFloat("Error threshold", form.Get("threshold.error"))
	var e float64
	if err != nil {
		errs = append(errs, err)
	} else if ok {
		e = f
	} else {
		errs = append(errs, errors.New("Error threshold should not be empty"))
	}
	f, err, ok = checkFloat("Critical threshold", form.Get("threshold.critical"))
	var c float64
	if err != nil {
		errs = append(errs, err)
	} else if ok {
		c = f
	} else {
		errs = append(errs, errors.New("Critical threshold should not be empty"))
	}

	if len(errs) > 0 {
		return "", errs
	}

	gt := newGraphiteThresholdSettings(graphite, metric, triggerIf, w, e, c, hold, alertFreq, checkFreq)

	b, err := json.Marshal(gt)
	if err != nil {
		fmt.Printf("Error with json: %s\n", err.Error())
		errs = append(errs, err)
	}
	s := string(b[:])
	return s, errs
}

func checkInt(t, s string) (uint, error, bool) {
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, errors.New(t + ": " + s + " should be a valid integer"), false
		}
		if i <= 0 {
			return 0, errors.New(t + ": " + s + " should be greater than 0"), false
		}

		return uint(i), nil, true
	}
	return 0, nil, false
}

func checkFloat(t, s string) (float64, error, bool) {
	if s != "" {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, errors.New(t + ": " + s + " should be a valid float"), false
		}
		if f <= 0 {
			return 0, errors.New(t + ": " + s + " should be greater than 0"), false
		}
		return f, nil, true
	}
	return 0.0, nil, false
}

func (gt GraphiteThreshold) AlertFrequency() uint {
	return gt.alertFrequency
}

func (gt *GraphiteThreshold) Check() (map[string]revere.Reading, error) {
	time := gt.checkFrequency * gt.holdRequirement
	resp, err := http.Get(
		fmt.Sprintf(
			graphiteUrlFormat,
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
		if trhld, ok := gt.thresholds[revere.Warning]; ok {
			if compare(gt.triggerIf, trhld, target.Datapoints) {
				reading.State = revere.Warning
			}
		}
		if trhld, ok := gt.thresholds[revere.Error]; ok {
			if compare(gt.triggerIf, trhld, target.Datapoints) {
				reading.State = revere.Error
			}
		}
		if trhld, ok := gt.thresholds[revere.Critical]; ok {
			if compare(gt.triggerIf, trhld, target.Datapoints) {
				reading.State = revere.Critical
			}
		}

		reading.Details = graphiteReadingDetails{
			"Target " +
				target.Target +
				" has state " +
				reading.State.String() +
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
