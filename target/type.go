package target

import (
	"github.com/jmoiron/sqlx/types"

	"github.com/yext/revere/db"
)

// Type defines a common abstraction for types of targets.
type Type interface {
	ID() db.TargetType

	// New returns a new instance of a target of this type.
	New(config types.JSONText) (Target, error)

	// Alert sends alert a to the targets in toAlert.
	//
	// Targets of this type that are also applicable to the alerting
	// subprobe but which shouldn't be alerted (e.g., because they've been
	// alerted too recently or they shouldn't be triggered until a higher
	// level) are included in inactive.
	//
	// Alert returns errors encountered during the sending of alerts,
	// associated with the trigger IDs of the targets that failed with those
	// errors.
	Alert(a *Alert, toAlert map[db.TriggerID]Target, inactive []Target) []ErrorAndTriggerIDs
}

type ErrorAndTriggerIDs struct {
	Err error
	IDs []db.TriggerID
}
