package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
)

type Label struct {
	*revere.Label
	Monitors []*LabelMonitor
	Triggers []*Trigger
}

func (l *Label) Id() int64 {
	return int64(l.Label.LabelId)
}

func NewLabel(db *sql.DB, id revere.LabelID) (*Label, error) {
	l, err := revere.LoadLabel(db, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, fmt.Errorf("Label not found: %d", id)
	}

	return newLabelFromModel(db, l)
}

func newLabelFromModel(db *sql.DB, l *revere.Label) (*Label, error) {
	var err error
	viewmodel := new(Label)
	viewmodel.Label = l
	viewmodel.Triggers, err = NewTriggersFromLabelTriggers(l.Triggers)
	if err != nil {
		return nil, err
	}
	viewmodel.Monitors = NewLabelMonitors(l.Monitors)
	if err != nil {
		return nil, err
	}
	return viewmodel, nil
}

func newLabelsFromModels(db *sql.DB, rls []*revere.Label) []*Label {
	labels := make([]*Label, len(rls))
	for i, rl := range rls {
		labels[i] = new(Label)
		labels[i].Label = rl
	}
	return labels
}

func BlankLabel(db *sql.DB) (*Label, error) {
	var err error
	viewmodel := new(Label)
	viewmodel.Label = new(revere.Label)
	viewmodel.Triggers = []*Trigger{
		BlankTrigger(),
	}
	viewmodel.Monitors = BlankLabelMonitors()
	if err != nil {
		return nil, err
	}

	return viewmodel, nil
}

func AllLabels(db *sql.DB) ([]*Label, error) {
	rls, err := revere.LoadLabels(db)
	if err != nil {
		return nil, err
	}

	return newLabelsFromModels(db, rls), nil
}
