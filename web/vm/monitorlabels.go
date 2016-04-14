package vm

import (
	"database/sql"

	"github.com/yext/revere"
)

type MonitorLabel struct {
	*revere.MonitorLabel
}

func (ml *MonitorLabel) Id() int64 {
	return int64(ml.MonitorLabel.LabelId)
}

func (ml *MonitorLabel) Name() string {
	return ml.MonitorLabel.Name
}

func (ml *MonitorLabel) Description() string {
	return ml.MonitorLabel.Description
}

func (ml *MonitorLabel) Subprobe() string {
	return ml.MonitorLabel.Subprobe
}

func NewMonitorLabels(monitorLabels []*revere.MonitorLabel) []*MonitorLabel {
	viewmodels := make([]*MonitorLabel, len(monitorLabels))
	for i, ml := range monitorLabels {
		viewmodels[i] = &MonitorLabel{ml}
	}
	return viewmodels
}

func BlankMonitorLabels() []*MonitorLabel {
	return []*MonitorLabel{}
}

func allMonitorLabels(db *sql.DB, mIds []revere.MonitorID) (map[revere.MonitorID][]*MonitorLabel, error) {
	labelsByMonitorId, err := revere.BatchLoadMonitorLabels(db, mIds)
	if err != nil {
		return nil, err
	}

	mls := make(map[revere.MonitorID][]*MonitorLabel)
	for mId, labels := range labelsByMonitorId {
		mls[mId] = NewMonitorLabels(labels)
	}
	return mls, nil
}
