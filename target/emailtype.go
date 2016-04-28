package target

import (
	"fmt"

	"github.com/jmoiron/sqlx/types"

	"github.com/yext/revere/db"
)

type emailType struct{}

func (_ emailType) ID() db.TargetType {
	return 1
}

func (_ emailType) New(config types.JSONText) (Target, error) {
	return newEmail(config)
}

func (_ emailType) Alert(a *Alert, toAlert map[db.TriggerID]Target, inactive []Target) []ErrorAndTriggerIDs {
	// TODO(eefi): Implement.
	fmt.Printf("Planning to email for %s/%s to %d targets.\n", a.MonitorName, a.SubprobeName, len(toAlert))
	return nil
}
