package target

import (
	"github.com/jmoiron/sqlx/types"

	"github.com/yext/revere/db"
)

type emailType struct{}

func (_ emailType) New(config types.JSONText) (Target, error) {
	return newEmail(config)
}

func (_ emailType) ID() db.TargetType {
	return 1
}
