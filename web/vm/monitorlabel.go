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

	mls := newMonitorLabelsFromDB(monitorLabels)
	return mls, nil
}

func newMonitorLabelFromDB(monitorLabel *db.MonitorLabel) *MonitorLabel {
	return &MonitorLabel{
		Label:     newLabelFromDB(monitorLabel.Label),
		MonitorID: monitorLabel.MonitorID,
		Subprobes: monitorLabel.Subprobes,
	}
}

func newMonitorLabelsFromDB(monitorLabels []*db.MonitorLabel) []*MonitorLabel {
	mls := make([]*MonitorLabel, len(monitorLabels))
	for i, monitorLabel := range monitorLabels {
		mls[i] = newMonitorLabelFromDB(monitorLabel)
	}
	return mls
}

func blankMonitorLabels() []*MonitorLabel {
	return []*MonitorLabel{}
}

func (ml *MonitorLabel) IsCreate() bool {
	return ml.Create
}

func (ml *MonitorLabel) IsDelete() bool {
	return ml.Delete
}

func (ml *MonitorLabel) validate(DB *db.DB) (errs []string) {
	// TODO(fchen) probably will want to start doing validation in transaction
	// also need to verify MonitorID is consistent with parent monitor
	if !DB.IsExistingLabel(db.LabelID(ml.Label.Id())) {
		errs = append(errs, fmt.Sprintf("Invalid label: %d", db.LabelID(ml.Label.Id())))
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
		return nil, errors.Trace(err)
	}

	mls := make(map[db.MonitorID][]*MonitorLabel)
	for mId, labels := range labelsByMonitorId {
		mls[mId] = newMonitorLabelsFromDB(labels)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return mls, nil
}

func (ml *MonitorLabel) setMonitorID(id db.MonitorID) {
	ml.MonitorID = id
}
