package vm

import (
	"database/sql"
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type MonitorTrigger struct {
	Trigger   *Trigger
	MonitorID db.MonitorID
	Subprobes string
	Delete    bool
}

func newMonitorTriggers(tx *db.Tx, id db.MonitorID) ([]MonitorTrigger, error) {
	monitorTriggers := db.LoadMonitorTriggers(tx, id)

	mts := make([]MonitorTrigger, len(monitorTriggers))
	var err error
	for i, monitorTrigger := range monitorTriggers {
		mts[i].Trigger, err = newTriggerFromModel(monitorTrigger.Trigger)
		if err != nil {
			return nil, errors.Trace(err)
		}
		mts[i].MonitorID = monitorTrigger.MonitorID
		mts[i].Subprobes = monitorTrigger.Subprobes
	}

	return mts
}

func blankMonitorTriggers() []MonitorTrigger {
	return []MonitorTrigger{}
}

func (mt *MonitorTrigger) Del() {
	return mt.Delete
}

func (mt *MonitorTrigger) Validate() {
	//TODO(fchen): maybe validate that MonitorID exists and is legitimate
	if !db.IsExistingMonitor(mt.MonitorID) {
		errs = append(errs, fmt.Sprintf("Invalid monitor: %d", mt.MonitorID))
	}
	if err := validateSubprobe(mt.Subprobe); err != nil {
		errs = append(errs, err.Error())
	}
	errs = append(errs, mt.Trigger.Validate()...)
	return
}

func (mt *MonitorTrigger) save(tx *sql.Tx, id db.MonitorID) error {
	monitorTrigger := &db.MonitorTrigger{mt.Trigger.toModelTrigger(), mt.Subprobe}
	var err error
	if isCreate(mt) {
		id, err := monitorTrigger.create(tx, id)
		mt.Trigger.setId(id)
	} else if isDelete(mt) {
		err = monitorTrigger.delete(tx, id)
	} else {
		err = monitorTrigger.update(tx, id)
	}

	return errors.Trace(err)
}
