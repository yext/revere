// Package probe implements the ways Revere can measure the outside world.
package probe

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
)

// Probe defines a common abstraction for all probes.
//
// TODO(eefi): Pull in the probe concept explanation text from the Revere
// product overview design doc. Also document all the individual methods.
//
// TODO(eefi): Incorporate the web side and also document that.
type Probe interface {
	Start()
	Stop()
}

// Readings is a group of readings from all a probe's subprobes taken at a
// specific time.
type Readings struct {
	Recorded time.Time
	Readings []Reading
}

// Reading is a specific reading from a particular subprobe.
type Reading struct {
	Subprobe string
	State    state.State
	Details  Details
}

// Details encodes probe-type-specific details from a reading.
type Details interface {
	// TODO(eefi): Fill in.
}

// New makes a Probe of the given type and settings. The Probe will send
// its readings to the provided channel.
func New(typeID db.ProbeType, config types.JSONText, readingsSink chan<- *Readings) (Probe, error) {
	// TODO(eefi): Implement Type dictionary system.
	if typeID != 1 {
		return nil, errors.Errorf("unknown probe type %d", typeID)
	}

	return graphiteThresholdType{}.New(config, readingsSink)
}
