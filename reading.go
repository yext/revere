package revere

// A Reading captures a snapshot of the state of a service being monitored.
type Reading struct {
	State   State
	Details Details
}

const (
	CheckFrequency = 1
)

// A State represents the health of a service being monitored.
type State int

const (
	// Things are operating as expected.
	Normal State = iota

	// There is something odd happening. People actively working on this
	// service or checking up on it should know about this, but others do
	// not need to know about it.
	Warning

	// There is a problem that needs to be addressed. People should be
	// alerted so that they can take steps to resolve the problem in the
	// near future.
	Error

	// There is a serious problem that needs immediate attention. The
	// service is in a state that warrants waking people up if necessary so
	// that the problem can be resolved as soon as possible.
	Critical

	// There has been an error in the health checking mechanism for this
	// service.
	Unknown State = -1
)

// Details provides a standard interface for manipulating Probe-specific
// details from a Reading.
type Details interface {
	// Text returns a human-readable plain-text description of this
	// Details's data suitable for inclusion in emitted Alerts.
	Text() string
}
