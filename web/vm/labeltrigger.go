package vm

import (
	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type LabelTrigger struct {
	Trigger *Trigger
	LabelID db.LabelID
}

func newLabelTriggers(tx *db.Tx, id db.LabelID) ([]*LabelTrigger, error) {
	var err error
	labelTriggers, err := tx.LoadTriggersForLabel(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	lts := make([]*LabelTrigger, len(labelTriggers))
	for i, labelTrigger := range labelTriggers {
		t, err := newTriggerFromModel(labelTrigger.Trigger)
		if err != nil {
			return nil, errors.Trace(err)
		}
		lts[i] = &LabelTrigger{
			Trigger: t,
			LabelID: labelTrigger.LabelID,
		}
	}

	return lts, nil
}

func BlankLabelTrigger() *LabelTrigger {
	return &LabelTrigger{
		Trigger: BlankTrigger(),
	}
}

func blankLabelTriggers() []*LabelTrigger {
	return []*LabelTrigger{}
}

func (lt *LabelTrigger) Id() int64 {
	return lt.Trigger.Id()
}

func (lt *LabelTrigger) IsCreate() bool {
	return lt.Id() == 0
}

func (lt *LabelTrigger) IsDelete() bool {
	return lt.Trigger.Delete
}

func (lt *LabelTrigger) validate(db *db.DB) (errs []string) {
	return append(errs, lt.Trigger.validate()...)
}

func (lt *LabelTrigger) save(tx *db.Tx) error {
	var err error
	trigger, err := lt.Trigger.toDBTrigger()
	if err != nil {
		return errors.Trace(err)
	}
	labelTrigger := db.LabelTrigger{
		LabelID: lt.LabelID,
		Trigger: trigger,
	}
	if isCreate(lt) {
		var id db.TriggerID
		id, err = tx.CreateLabelTrigger(labelTrigger)
		lt.Trigger.setId(id)
	} else if isDelete(lt) {
		err = tx.DeleteLabelTrigger(labelTrigger.TriggerID)
	} else {
		err = tx.UpdateLabelTrigger(labelTrigger)
	}

	return errors.Trace(err)
}

func (lt *LabelTrigger) setLabelID(id db.LabelID) {
	lt.LabelID = id
}
