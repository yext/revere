package vm

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type LabelMonitor struct {
	Monitor   *Monitor
	LabelID   db.LabelID
	Subprobes string
	Delete    bool
}

func (lm *LabelMonitor) Id() int64 {
	return int64(lm.LabelMonitor.MonitorId)
}

func newLabelMonitors(tx *db.Tx, id db.LabelID) ([]LabelMonitor, error) {
	labelMonitors, err := tx.LoadMonitorsForLabel(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	lms := make([]LabelMonitor, len(labelMonitors))
	for i, labelMonitor := range labelMonitors {
		lms[i].Monitor, err = newMonitorFromModel(labelMonitor.Monitor)
		if err != nil {
			return nil, errors.Trace(err)
		}
		lms[i].LabelID = labelMonitor.LabelID
		lms[i].Subprobes = labelMonitor.Subprobes
	}

	return lms
}

func blankLabelMonitors() []LabelMonitor {
	return []LabelMonitor{}
}

func (lm *LabelMonitor) Del() {
	return lm.Delete
}

func (lm *LabelMonitor) validate(db *db.DB) (errs []string) {
	if !db.IsExistingLabel(lm.LabelID) {
		errs = append(errs, fmt.Sprintf("Invalid label: %d", lm.LabelID))
	}
	if err := validateSubprobeRegex(lm.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}
	return
}

func (lm *LabelMonitor) save(tx *db.Tx, id db.LabelID) error {
	monitor, err := lm.Monitor.toModelMonitor()
	if err != nil {
		return errors.Trace(err)
	}
	labelMonitor := &db.LabelMonitor{
		LabelID:   lm.LabelID,
		Subprobes: lm.Subprobes,
		Monitor:   monitor,
	}
	if isCreate(lm) {
		err = tx.CreateLabelMonitor(labelMonitor)
	} else if isDelete(lm) {
		err = tx.DeleteLabelMonitor(labelMonitor)
	} else {
		err = tx.UpdateLabelMonitor(labelMonitor)
	}

	return errors.Trace(err)
}
