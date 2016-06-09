package vm

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type Label struct {
	LabelID     db.LabelID
	Name        string
	Description string
	Triggers    []*LabelTrigger
	Monitors    []*LabelMonitor
}

func (l *Label) Id() int64 {
	return int64(l.LabelID)
}

func NewLabel(tx *db.Tx, id db.LabelID) (*Label, error) {
	label, err := tx.LoadLabel(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if label == nil {
		return nil, errors.Errorf("Label not found: %d", id)
	}

	l := newLabelFromDB(label)

	err := l.loadComponents(tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return l, nil
}

func newLabelFromDB(label *db.Label) *Label {
	return &Label{
		LabelID:     label.LabelID,
		Name:        label.Name,
		Description: label.Description,
		Triggers:    nil,
		Monitors:    nil,
	}
}

func newLabelsFromDB(labels []*db.Label) []*Label {
	ls := make([]*Label, len(labels))
	for i, label := range labels {
		ls[i] = newLabelFromDB(label)
	}
	return ls
}

func BlankLabel() (*Monitor, error) {
	var err error
	l := &Label{}
	l.Triggers = blankLabelTriggers()
	l.Monitors = blankLabelMonitors()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return l, nil
}

func AllLabels(tx *db.Tx) ([]*Label, error) {
	ls, err := tx.LoadLabels()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return newLabelsFromDB(ls), nil
}

func (l *Label) loadComponents(tx *db.Tx) error {
	var err error
	l.Triggers, err = newLabelTriggers(tx, l.LabelID)
	if err != nil {
		return errors.Trace(err)
	}
	l.Monitors, err = NewLabelMonitors(tx, l.LabelID)
	return errors.Trace(err)
}

func (l *Label) Validate(db *db.DB) (errs []string) {
	if l.Name == "" {
		errs = append(errs, fmt.Sprintf("Label name is required"))
	}

	for _, lt := range l.Triggers {
		errs = append(errs, lt.validate()...)
	}

	for _, lm := range l.Monitors {
		errs = append(errs, lm.validate(db)...)
	}
	return
}

func (l *Label) Create() bool {
	return l.Id() == 0
}

func (l *Label) Save(tx *db.DB) error {
	label := l.toDBLabel()

	var err error
	if isCreate(l) {
		l.LabelID, err = db.CreateLabel(label)
	} else {
		err = db.UpdateLabel(label)
	}
	if err != nil {
		return errors.Trace(err)
	}

	for _, t := range l.Triggers {
		err = t.save(tx)
		if err != nil {
			return err
		}
	}

	for _, m := range l.Monitors {
		err = m.save(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Label) toDBLabel() *db.Label {
	return &db.Label{
		LabelID:     l.LabelID,
		Name:        l.Name,
		Description: l.Description,
	}
}
