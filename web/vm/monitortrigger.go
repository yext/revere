package vm

import (
	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type MonitorTrigger struct {
	Trigger   *Trigger
	MonitorID db.MonitorID
	Subprobes string
}

func newMonitorTriggers(tx *db.Tx, id db.MonitorID) ([]*MonitorTrigger, error) {
	monitorTriggers, err := tx.LoadTriggersForMonitor(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	mts := make([]*MonitorTrigger, len(monitorTriggers))
	for i, monitorTrigger := range monitorTriggers {
		t, err := newTriggerFromModel(monitorTrigger.Trigger)
		if err != nil {
			return nil, errors.Trace(err)
		}
		mts[i] = &MonitorTrigger{
			Trigger:   t,
			MonitorID: monitorTrigger.MonitorID,
			Subprobes: monitorTrigger.Subprobes,
		}
	}

	return mts, nil
}

func BlankMonitorTrigger() *MonitorTrigger {
	return &MonitorTrigger{
		Trigger:   BlankTrigger(),
		Subprobes: "",
	}
}

func blankMonitorTriggers() []*MonitorTrigger {
	return []*MonitorTrigger{}
}

func (mt *MonitorTrigger) validate(db *db.DB) (errs []string) {
	if err := validateSubprobeRegex(mt.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}
	errs = append(errs, mt.Trigger.validate()...)
	return
}

func (mt *MonitorTrigger) Id() int64 {
	return mt.Trigger.Id()
}

func (mt *MonitorTrigger) IsCreate() bool {
	return mt.Id() == 0
}

func (mt *MonitorTrigger) IsDelete() bool {
	return mt.Trigger.Delete
}

func (mt *MonitorTrigger) save(tx *db.Tx) error {
	trigger, err := mt.Trigger.toDBTrigger()
	if err != nil {
		return errors.Trace(err)
	}
	monitorTrigger := db.MonitorTrigger{
		MonitorID: mt.MonitorID,
		Subprobes: mt.Subprobes,
		Trigger:   trigger,
	}
	if isCreate(mt) {
		var id db.TriggerID
		id, err = tx.CreateMonitorTrigger(monitorTrigger)
		mt.Trigger.setId(id)
	} else if isDelete(mt) {
		err = tx.DeleteMonitorTrigger(monitorTrigger.TriggerID)
	} else {
		err = tx.UpdateMonitorTrigger(monitorTrigger)
	}

	return errors.Trace(err)
}

func (mt *MonitorTrigger) setMonitorID(id db.MonitorID) {
	mt.MonitorID = id
}
