package vm

import (
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
	monitorTriggers, err := tx.LoadTriggersForMonitor(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	mts := make([]MonitorTrigger, len(monitorTriggers))
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

func (mt *MonitorTrigger) validate(db *db.DB) (errs []string) {
	if !db.IsExistingMonitor(mt.MonitorID) {
		errs = append(errs, fmt.Sprintf("Invalid monitor: %d", mt.MonitorID))
	}
	if err := validateSubprobeRegex(mt.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}
	errs = append(errs, mt.Trigger.validate()...)
	return
}

func (mt *MonitorTrigger) save(tx *db.Tx, id db.MonitorID) error {
	monitorTrigger := &db.MonitorTrigger{
		MonitorID: mt.MonitorID,
		Subprobes: mt.Subprobes,
		Trigger:   mt.Trigger.toModelTrigger(),
	}
	var err error
	if isCreate(mt) {
		id, err := tx.CreateMonitorTrigger(monitorTrigger)
		mt.Trigger.setId(id)
	} else if isDelete(mt) {
		err = tx.DeleteMonitorTrigger(monitorTrigger.TriggerID)
	} else {
		err = tx.UpdateMonitorTrigger(monitorTrigger)
	}

	return errors.Trace(err)
}
