package probe

import (
	"github.com/jmoiron/sqlx/types"
)

type graphiteThresholdType struct{}

func (t graphiteThresholdType) New(config types.JSONText, readingsSink chan<- Readings) (Probe, error) {
	return newGraphiteThreshold(config, readingsSink)
}
