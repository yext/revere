package probe

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"
)

// GraphiteThreshold implements a probe that assigns states based on whether a
// Graphite metric is above or below various constant values.
type GraphiteThreshold struct {
	*GraphiteThresholdDBModel

	readingsSink chan<- *Readings
}

func newGraphiteThreshold(config types.JSONText, readingsSink chan<- *Readings) (Probe, error) {
	var json GraphiteThresholdDBModel
	err := config.Unmarshal(&json)
	if err != nil {
		return nil, errors.Maskf(err, "deserialize probe config")
	}

	return &GraphiteThreshold{
		GraphiteThresholdDBModel: &json,
		readingsSink:             readingsSink,
	}, nil
}

func (gt *GraphiteThreshold) Start() {
}

func (gt *GraphiteThreshold) Stop() {
}
