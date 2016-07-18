package vm

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type LabelMonitor struct {
	// TODO(fchen): Organize struct fields between db-only fields, front-end-only fields, and shared fields
	Monitor   *Monitor
	LabelID   db.LabelID
	Subprobes string
	Create    bool
	Delete    bool
}

func (lm *LabelMonitor) Id() int64 {
	return int64(lm.Monitor.MonitorID)
}

func newLabelMonitors(tx *db.Tx, id db.LabelID) ([]*LabelMonitor, error) {
	labelMonitors, err := tx.LoadMonitorsForLabel(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	lms := make([]*LabelMonitor, len(labelMonitors))
	for i, labelMonitor := range labelMonitors {
		m, err := newMonitorFromDB(labelMonitor.Monitor, tx)
		if err != nil {
			return nil, errors.Trace(err)
		}
		lms[i] = &LabelMonitor{
			Monitor:   m,
			LabelID:   labelMonitor.LabelID,
			Subprobes: labelMonitor.Subprobes,
		}
	}

	return lms, nil
}

func blankLabelMonitors() []*LabelMonitor {
	return []*LabelMonitor{}
}

func (lm *LabelMonitor) IsCreate() bool {
	return lm.Create
}

func (lm *LabelMonitor) IsDelete() bool {
	return lm.Delete
}

func (lm *LabelMonitor) validate(DB *db.DB) (errs []string) {
	// TODO(fchen) probably will want to start doing validation in transaction
	// also need to verify LabelID is consistent with parent label
	if !DB.IsExistingMonitor(db.MonitorID(lm.Monitor.Id())) {
		errs = append(errs, fmt.Sprintf("Invalid monitor: %d", db.MonitorID(lm.Monitor.Id())))
	}
	if err := validateSubprobeRegex(lm.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}
	return
}

func (lm *LabelMonitor) save(tx *db.Tx) error {
	labelMonitor := db.LabelMonitor{
		LabelID:   lm.LabelID,
		Subprobes: lm.Subprobes,
		Monitor: &db.Monitor{
			MonitorID: lm.Monitor.MonitorID,
		},
	}

	var err error
	if isCreate(lm) {
		err = tx.CreateLabelMonitor(labelMonitor)
	} else if isDelete(lm) {
		err = tx.DeleteLabelMonitor(labelMonitor)
	} else {
		err = tx.UpdateLabelMonitor(labelMonitor)
	}

	return errors.Trace(err)
}

func (lm *LabelMonitor) setLabelID(id db.LabelID) {
	lm.LabelID = id
}
