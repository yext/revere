package daemon

import (
	"time"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
)

type subprobe struct {
	id      db.SubprobeID
	monitor *monitor
	name    string

	lastReading  time.Time
	state        state.State
	enteredState time.Time
	lastNormal   time.Time

	triggerSets map[db.TargetType]*sameTypeTriggerSet
}

func newSubprobe(name string, status db.SubprobeStatus, monitor *monitor) *subprobe {
	triggerSets := make(map[db.TargetType]*sameTypeTriggerSet)
	for _, monitorTrigger := range monitor.triggers {
		if monitorTrigger.subprobes.MatchString(name) {
			trigger := monitorTrigger.trigger
			targetType := trigger.target.Type().ID()

			triggerSet := triggerSets[targetType]
			if triggerSet == nil {
				triggerSet = &sameTypeTriggerSet{}
				triggerSets[targetType] = triggerSet
			}

			triggerSet.add(trigger)
		}
	}

	return &subprobe{
		id:           status.SubprobeID,
		monitor:      monitor,
		name:         name,
		lastReading:  status.Recorded,
		state:        status.State,
		enteredState: status.EnteredState,
		lastNormal:   status.LastNormal,
		triggerSets:  triggerSets,
	}
}
