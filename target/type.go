package target

import (
	"github.com/yext/revere/db"
)

// Type defines a common abstraction for types of targets.
type Type interface {
	ID() db.TargetType
}
