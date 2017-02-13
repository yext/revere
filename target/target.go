// Package target implements the ways Revere can alert people to problems.
package target

import (
	"fmt"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
)

var (
	daemonTargetTypes = make(map[db.TargetType]Type)
)

// Target defines a common abstraction for individual targets.
type Target interface {
	Type() Type
}

// New makes a Target of the given type and settings.
func New(typeID db.TargetType, config types.JSONText) (Target, error) {
	if targetType, found := daemonTargetTypes[typeID]; found {
		return targetType.New(config)
	}
	return nil, errors.Errorf("unknown target type %d", typeID)
}

// registerTargetType registers a target type onto a type dictionary
func registerTargetType(t Type) {
	if _, exists := daemonTargetTypes[t.ID()]; !exists {
		daemonTargetTypes[t.ID()] = t
	} else {
		panic(fmt.Sprintf("A target type with id %d already exists", t.ID()))
	}
}
