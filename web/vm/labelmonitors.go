package vm

import (
	"database/sql"

	"github.com/yext/revere"
)

type LabelMonitors struct {
	Components    []*revere.LabelMonitor
	All           []*revere.Monitor
	ComponentName string
}

func NewLabelMonitors(db *sql.DB, labelMonitors []*revere.LabelMonitor) (*LabelMonitors, error) {
	viewmodel, err := loadLabelMonitors(db)
	if err != nil {
		return nil, err
	}
	viewmodel.Components = labelMonitors
	return viewmodel, nil
}

func BlankLabelMonitors(db *sql.DB) (*LabelMonitors, error) {
	return loadLabelMonitors(db)
}

func loadLabelMonitors(db *sql.DB) (*LabelMonitors, error) {
	var err error
	viewmodel := new(LabelMonitors)
	viewmodel.ComponentName = "monitors"
	viewmodel.All, err = revere.LoadMonitors(db)
	if err != nil {
		return nil, err
	}
	return viewmodel, nil
}
