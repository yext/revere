package probe

import (
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/state"
)

// GraphiteThreshold implements a probe that assigns states based on whether a
// Graphite metric is above or below various constant values.
type GraphiteThreshold struct {
	*Polling
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

	return &gt, nil
}

func (gt *GraphiteThreshold) Check() []Reading {
	// TODO(eefi): Implement.

	s := state.Normal
	if rand.Intn(4) == 0 {
		s = state.Error
	}

	return []Reading{
		{
			Subprobe: "_",
			State:    s,
			Recorded: time.Now(),
		},
	}
}
