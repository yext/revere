package vm

import (
	"html/template"

	"github.com/yext/revere"
)

type Label struct {
	*revere.Label
	Monitors    []*revere.LabelMonitor
	Triggers    []*revere.LabelTrigger
	AllMonitors []*revere.Monitor
}

func NewLabel(l *revere.Label, allMonitors []*revere.Monitor) (*Label, error) {
	viewmodel := new(Label)
	viewmodel.Label = l
	viewmodel.Triggers = l.Triggers
	viewmodel.Monitors = l.Monitors
	viewmodel.AllMonitors = allMonitors
	// TODO(psingh): Take out error if not needed
	return viewmodel, nil
}

func BlankLabel() (*Label, error) {
	viewmodel := new(Label)
	viewmodel.Label = new(revere.Label)
	viewmodel.Triggers = []*revere.LabelTrigger{
		&revere.LabelTrigger{
			Trigger: revere.Trigger{
				TargetTemplate: template.HTML("PLACEHOLDER DELETE"), //TODO(fchen): code cleanup[targets.DefaultTargetTemplate()]
			},
		},
	}

	// TODO(psingh): Add monitor related stuff
	return viewmodel, nil
}
