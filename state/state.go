// Package state defines the set of states monitored services can be in in
// Revere's view.
package state

import (
	"fmt"

	"github.com/juju/errors"
)

// State is a state that Revere can consider a monitored service to be in. The
// available states are defined as constants in this package.
type State int8

// The states that Revere can consider a monitored service to be in. Although
// not required, it is highly recommended that people setting up monitors follow
// the suggested semantics documented here.
const (
	// Normal means things are operating as expected.
	Normal State = iota * 10

	// Warning means there is something odd happening. People actively
	// working on a service or checking up on it should know about the
	// situation, but others do not need to know about it.
	//
	// A service in Warning state shows up on Revere's dashboard of all
	// unhealthy services, and the times when a service slips into and out
	// of Warning state are recorded by Revere in a service's health
	// history, but entering Warning state generally does not warrant
	// triggering emails sent to a wide audience.
	Warning

	// Unknown means there has been an error in the health checking
	// mechanism for a service.
	//
	// For example, for a system monitored via a Graphite threshold probe,
	// if Graphite cannot be reached, or if the provided Graphite target
	// expression results in a Graphite error, then the service is in an
	// unknown state.
	Unknown

	// Error means there is a problem that needs to be addressed. People
	// should be alerted so that they can take steps to resolve the problem
	// in the near future.
	//
	// Entering Error state usually warrants emailing the mailing list of
	// the team that owns a service.
	Error

	// Critical means there is a serious problem that needs immediate
	// attention. A service has entered a state that warrants waking people
	// up if necessary so that the problem can be resolved as quickly as
	// possible.
	Critical
)

func (s State) String() string {
	switch s {
	case Normal:
		return "Normal"
	case Warning:
		return "Warning"
	case Unknown:
		return "Unknown"
	case Error:
		return "ERROR"
	case Critical:
		return "CRITICAL"
	default:
		return fmt.Sprintf("Invalid(%d)", s)
	}
}

func FromString(s string) (State, err) {
	switch s {
	case Normal.String():
		return Normal, nil
	case Warning.String():
		return Warning, nil
	case Unknown.String():
		return Unknown, nil
	case Error.String():
		return Error, nil
	case Critical.String():
		return Critical, nil
	default:
		return nil, errors.Errorf("invalid state %s", s)
	}
}
