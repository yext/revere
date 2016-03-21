package vm

import (
	"database/sql"

	"github.com/yext/revere"
)

type MonitorLabels struct {
	Components    []*revere.MonitorLabel
	All           []*revere.Label
	ComponentName string
}

func NewMonitorLabels(db *sql.DB, monitorLabel []*revere.MonitorLabel) (*MonitorLabels, error) {
	viewmodel, err := loadMonitorLabels(db)
	if err != nil {
		return nil, err
	}
	viewmodel.Components = monitorLabel
	return viewmodel, nil
}

func BlankMonitorLabels(db *sql.DB) (*MonitorLabels, error) {
	return loadMonitorLabels(db)
}

func loadMonitorLabels(db *sql.DB) (*MonitorLabels, error) {
	var err error
	viewmodel := new(MonitorLabels)
	viewmodel.ComponentName = "labels"
	viewmodel.All, err = revere.LoadLabels(db)
	if err != nil {
		return nil, err
	}
	return viewmodel, nil
}
