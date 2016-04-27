// Package target implements the ways Revere can alert people to problems.
package target

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
)

// Target defines a common abstraction for individual targets.
type Target interface {
	Type() Type
}

// New makes a Target of the given type and settings.
func New(typeID db.TargetType, config types.JSONText) (Target, error) {
	// TODO(eefi): Implement Type dictionary system.
	if typeID != 1 {
		return nil, errors.Errorf("unknown target type %d", typeID)
	}

	return emailType{}.New(config)
}
