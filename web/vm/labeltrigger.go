package vm

import (
	"database/sql"
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type LabelTrigger struct {
	Trigger *Trigger
	LabelID db.LabelID
	Delete  bool
}

func NewLabelTriggers(db *sql.DB, id db.LabelID) ([]LabelTrigger, error) {
	labelTriggers := db.LoadLabelTriggers(db, id)

	lts := make([]LabelTrigger, len(labelTriggers))
	var err error
	for i, labelTrigger := range labelTriggers {
		lts[i].Trigger, err = newTriggerFromModel(labelTrigger.Trigger)
		if err != nil {
			return nil, errors.Trace(err)
		}
		lts[i].LabelID = labelTrigger.LabelID
	}

	return lts
}

func BlankLabelTriggers() []LabelTrigger {
	return []LabelTrigger{}
}

func (lt *LabelTrigger) Del() {
	return lt.Delete
}

func (lt *LabelTrigger) validate(db *db.DB) (errs []string) {
	if !db.IsExistingLabel(lt.LabelID) {
		errs = append(errs, fmt.Sprintf("Invalid label: %d", lt.LabelID))
	}

	return append(errs, lt.Trigger.Validate()...)
}

func (lt *LabelTrigger) save(tx *sql.Tx, id db.LabelID) error {
	labelTrigger := &db.LabelTrigger{lt.Trigger.toModelTrigger(), lt.Subprobe}
	var err error
	if isCreate(lt) {
		id, err := labelTrigger.create(tx, id)
		lt.Trigger.setId(id)
	} else if isDelete(lt) {
		err = labelTrigger.delete(tx, id)
	} else {
		err = labelTrigger.update(tx, id)
	}

	return errors.Trace(err)
}
