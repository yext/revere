package web

import (
	"errors"

	"github.com/yext/revere"
)

// Provides a way to pass multiple values to a subtemplate
// Taken from: http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters
func dict(args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("Invalid dict call")
	}
	dict := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return nil, errors.New("Dict keys must be strings")
		}
		dict[key] = args[i+1]
	}
	return dict, nil
}

// Lookup threshold values, can't look up with consts in templates
func lookupThreshold(thresholds map[revere.State]float64, state string) (float64, error) {
	if state == "Warning" {
		return thresholds[revere.Warning], nil
	} else if state == "Error" {
		return thresholds[revere.Error], nil
	} else if state == "Critical" {
		return thresholds[revere.Critical], nil
	}
	return float64(revere.Unknown), nil
}
