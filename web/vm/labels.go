package vm

import (
	"github.com/yext/revere"
)

type Label struct {
	*revere.Label
	Monitors    []*revere.LabelMonitor
	Triggers    []*Trigger
	AllMonitors []*revere.Monitor
}

func NewLabel(l *revere.Label, allMonitors []*revere.Monitor) (*Label, error) {
	viewmodel := new(Label)

	var err error
	viewmodel.Label = l
	viewmodel.Triggers, err = NewTriggersFromLabelTriggers(l.Triggers)
	if err != nil {
		return nil, err
	}

	viewmodel.Monitors = l.Monitors
	viewmodel.AllMonitors = allMonitors

	return viewmodel, nil
}

func BlankLabel() (*Label, error) {
	viewmodel := new(Label)
	viewmodel.Label = new(revere.Label)
	viewmodel.Triggers = []*Trigger{
		BlankTrigger(),
	}

	// TODO(psingh): Add monitor related stuff
	return viewmodel, nil
}
