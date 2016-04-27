package probe

import (
	"github.com/jmoiron/sqlx/types"
)

type graphiteThresholdType struct{}

func (_ graphiteThresholdType) New(config types.JSONText, readingsSink chan<- []Reading) (Probe, error) {
	return newGraphiteThreshold(config, readingsSink)
}
