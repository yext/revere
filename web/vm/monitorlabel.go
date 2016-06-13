package vm

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type MonitorLabel struct {
	Label     *Label
	MonitorID db.MonitorID
	Subprobes string
	Create    bool
	Delete    bool
}

func (ml *MonitorLabel) Id() int64 {
	return int64(ml.Label.LabelID)
}

func newMonitorLabels(tx *db.Tx, id db.MonitorID) ([]*MonitorLabel, error) {
	monitorLabels, err := tx.LoadLabelsForMonitor(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	mls := make([]*MonitorLabel, len(monitorLabels))
	for i, monitorLabel := range monitorLabels {
		mls[i].Label = newLabelFromDB(monitorLabel.Label)
		mls[i].MonitorID = monitorLabel.MonitorID
		mls[i].Subprobes = monitorLabel.Subprobes
	}

	return mls, nil
}

func blankMonitorLabels() []*MonitorLabel {
	return []*MonitorLabel{}
}

func (ml *MonitorLabel) Create() bool {
	return ml.Create
}

func (ml *MonitorLabel) Del() bool {
	return ml.Delete
}

func (ml *MonitorLabel) validate(db *db.DB) (errs []string) {
	if !db.IsExistingMonitor(ml.MonitorID) {
		errs = append(errs, fmt.Sprintf("Invalid monitor: %d", ml.MonitorID))
	}
	if err := validateSubprobeRegex(ml.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}
	return
}

func (ml *MonitorLabel) save(tx *db.Tx) error {
	monitorLabel := db.MonitorLabel{
		MonitorID: ml.MonitorID,
		Subprobes: ml.Subprobes,
		Label:     ml.Label.toDBLabel(),
	}
	var err error
	if isCreate(ml) {
		err = tx.CreateMonitorLabel(monitorLabel)
	} else if isDelete(ml) {
		err = tx.DeleteMonitorLabel(monitorLabel)
	} else {
		err = tx.UpdateMonitorLabel(monitorLabel)
	}

	return errors.Trace(err)
}

func allMonitorLabels(tx *db.Tx, mIds []db.MonitorID) (map[db.MonitorID][]*MonitorLabel, error) {
	labelsByMonitorId, err := tx.BatchLoadMonitorLabels(mIds)
	if err != nil {
		return nil, err
	}

	mls := make(map[db.MonitorID][]*MonitorLabel)
	for mId, labels := range labelsByMonitorId {
		mls[mId], err = newMonitorLabels(tx, mId)
		if err != nil {
			return nil, err
		}
	}
	return mls, nil
}
