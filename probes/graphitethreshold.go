package probes

import (
	"github.com/yext/revere"
)

// A GraphiteThreshold probe compares Graphite metrics to static values.
type GraphiteThreshold struct {
	graphite   string
	metric     string
	thresholds map[revere.State]float64
}

// NewGraphiteThreshold builds a probe based on serialized settings.
//
// TODO(eefi): Detail the serialization format.
func NewGraphiteThreshold(settings string) *GraphiteThreshold {
	// TODO(eefi): Implement me.
	return nil
}

func (*GraphiteThreshold) Check() map[string]revere.Reading {
	// TODO(eefi): Implement me.
	return make(map[string]revere.Reading)
}
