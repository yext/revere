package probe

import (
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

func newGraphiteThreshold(configJSON types.JSONText, readingsSink chan<- *Readings) (Probe, error) {
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

func (gt *GraphiteThreshold) Check() *Readings {
	// TODO(eefi): Implement.
	return &Readings{
		Recorded: time.Now(),
		Readings: []Reading{
			{
				Subprobe: "_",
				State:    state.Unknown,
			},
		},
	}
}
