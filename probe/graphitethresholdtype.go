package probe

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/yext/revere/db"
)

type graphiteThresholdType struct{}

// TODO: Figure out something better than passing the transaction all the way through
func (_ graphiteThresholdType) New(tx *db.Tx, config types.JSONText, readingsSink chan<- []Reading) (Probe, error) {
	return newGraphiteThreshold(tx, config, readingsSink)
}
