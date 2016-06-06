package db

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/state"
)

type TriggerID int32
type TargetType int16

type Trigger struct {
	TriggerID     TriggerID
	Level         state.State
	TriggerOnExit bool
	PeriodMilli   int32
	TargetType    TargetType
	Target        types.JSONText
}

func (tx *Tx) createTrigger(t *Trigger) (TriggerID, error) {
	q := `INSERT INTO pfx_triggers (level, triggeronexit, periodmilli, targettype, target)
	      VALUES (:level, :triggeronexit, :periodmilli, :targettype, :target)`
	result, err := tx.NamedExec(cq(tx, q), t)
	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}
	return TriggerID(id), nil
}

func (tx *Tx) updateTrigger(t *Trigger) error {
	q := `UPDATE pfx_triggers
	      SET level=:level,
	          triggeronexit=:triggeronexit,
	          periodmilli=:periodmilli,
	          targettype=:targettype,
	          target=:target
	      WHERE triggerid=:triggerid`
	_, err := tx.NamedExec(cq(tx, q), t)
	return errors.Trace(err)
}

func (tx *Tx) deleteTrigger(triggerID TriggerID) error {
	q := `DELETE FROM pfx_triggers WHERE triggerid=?`
	_, err := tx.Exec(cq(tx, q), triggerID)
	return errors.Trace(err)
}
