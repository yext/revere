package target

import (
	"time"

	"github.com/yext/revere/db"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/state"
)

// Alert contains all the information targets can use to construct their alerts.
type Alert struct {
	MonitorID    db.MonitorID
	MonitorName  string
	SubprobeID   db.SubprobeID
	SubprobeName string

	Description string
	Response    string

	OldState state.State
	NewState state.State

	Recorded     time.Time
	EnteredState time.Time
	LastNormal   time.Time

	Details probe.Details
}
