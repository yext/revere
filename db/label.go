package db

import (
	"database/sql"

	"github.com/juju/errors"
)

type LabelID int32

type Label struct {
	LabelID     LabelID
	Name        string
	Description string
}

type LabelTrigger struct {
	LabelID LabelID
	*Trigger
}

type LabelMonitor struct {
	LabelID   LabelID
	Subprobes string
	*Monitor
}

func (db *DB) LoadLabel(id LabelID) (*Label, error) {
	return loadLabel(db, id)
}

func (tx *Tx) LoadLabel(id LabelID) (*Label, error) {
	return loadLabel(tx, id)
}

func loadLabel(dt dbOrTx, id LabelID) (*Label, error) {
	dt = unsafe(dt)

	var m Label
	err := dt.Get(&m, cq(dt, "SELECT * FROM pfx_labels WHERE labelid = ?"), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &m, nil
}

func (db *DB) LoadLabels() ([]*Label, error) {
	return loadLabels(db)
}

func (tx *Tx) LoadLabels() ([]*Label, error) {
	return loadLabels(tx)
}

func loadLabels(dt dbOrTx) ([]*Label, error) {
	var labels []*Label
	if err := dt.Select(&labels, cq(dt, "SELECT * FROM pfx_labels ORDER BY name")); err != nil {
		return nil, errors.Trace(err)
	}
	return labels, nil
}

func (db *DB) IsExistingLabel(id LabelID) (exists bool) {
	if id == 0 {
		return false
	}

	q := `SELECT EXISTS (SELECT * FROM pfx_labels WHERE labelid = ?)`
	err := db.Get(&exists, cq(db, q), id)
	if err != nil {
		return false
	}
	return
}

func (db *DB) LoadTriggersForLabel(id LabelID) ([]LabelTrigger, error) {
	return loadTriggersForLabel(db, id)
}

func (tx *Tx) LoadTriggersForLabel(id LabelID) ([]LabelTrigger, error) {
	return loadTriggersForLabel(tx, id)
}

func loadTriggersForLabel(dt dbOrTx, id LabelID) ([]LabelTrigger, error) {
	dt = unsafe(dt)

	var lts []LabelTrigger
	q := `SELECT *
	      FROM pfx_label_triggers
	      JOIN pfx_triggers USING (triggerid)
	      WHERE pfx_label_triggers.labelid = ?`
	err := dt.Select(&lts, cq(dt, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return lts, nil
}

func (db *DB) LoadMonitorsForLabel(id LabelID) ([]LabelMonitor, error) {
	return loadMonitorsForLabel(db, id)
}

func (tx *Tx) LoadMonitorsForLabel(id LabelID) ([]LabelMonitor, error) {
	return loadMonitorsForLabel(tx, id)
}

func loadMonitorsForLabel(dt dbOrTx, id LabelID) ([]LabelMonitor, error) {
	dt = unsafe(dt)

	var lms []LabelMonitor
	q := `SELECT *
	      FROM pfx_labels_monitors
	      JOIN pfx_monitors USING (monitorid)
	      WHERE pfx_labels_monitors.labelid = ?
		  ORDER BY pfx_monitors.name`
	err := dt.Select(&lms, cq(dt, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return lms, nil
}

type LabelTriggerWithSubprobes struct {
	LabelTrigger
	Subprobes string
}

func (tx *Tx) LoadLabelTriggersForMonitor(id MonitorID) ([]LabelTriggerWithSubprobes, error) {
	var results []LabelTriggerWithSubprobes
	q := `SELECT labelid, triggerid, pfx_triggers.*, pfx_labels_monitors.subprobes
	      FROM pfx_labels_monitors
	      JOIN pfx_label_triggers USING (labelid)
	      JOIN pfx_triggers USING (triggerid)
	      WHERE pfx_labels_monitors.monitorid = ?`
	err := tx.Select(&results, cq(tx, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return results, nil
}

func (tx *Tx) CreateLabel(l *Label) (LabelID, error) {
	q := `INSERT INTO pfx_labels (name, description)
	      VALUES (:name, :description)`
	result, err := tx.NamedExec(cq(tx, q), l)
	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}
	return LabelID(id), nil
}

func (tx *Tx) UpdateLabel(l *Label) error {
	q := `UPDATE pfx_labels
	      SET name=:name, description=:description
	      WHERE labelid=:labelid`
	_, err := tx.NamedExec(cq(tx, q), l)
	return errors.Trace(err)
}

func (tx *Tx) CreateLabelMonitor(lm LabelMonitor) error {
	q := `INSERT INTO pfx_labels_monitors (labelid, monitorid, subprobes)
	      VALUES (:labelid, :monitorid, :subprobes)`
	_, err := tx.NamedExec(cq(tx, q), lm)
	return errors.Trace(err)
}

func (tx *Tx) UpdateLabelMonitor(lm LabelMonitor) error {
	q := `UPDATE pfx_labels_monitors
	      SET subprobes=:subprobes
	      WHERE labelid=:labelid AND monitorid=:monitorid`
	_, err := tx.NamedExec(cq(tx, q), lm)
	return errors.Trace(err)
}

func (tx *Tx) DeleteLabelMonitor(lm LabelMonitor) error {
	q := `DELETE FROM pfx_labels_monitors
	      WHERE labelid=:labelid AND monitorid=:monitorid`
	_, err := tx.NamedExec(cq(tx, q), lm)
	return errors.Trace(err)
}

func (tx *Tx) DeleteLabelTrigger(triggerID TriggerID) error {
	return tx.deleteTrigger(triggerID)
}

func (tx *Tx) CreateLabelTrigger(lt LabelTrigger) (TriggerID, error) {
	var err error
	lt.TriggerID, err = tx.createTrigger(lt.Trigger)
	if err != nil {
		return 0, errors.Trace(err)
	}

	q := `INSERT INTO pfx_label_triggers (labelid, triggerid)
	      VALUES (:labelid, :triggerid)`
	_, err = tx.NamedExec(cq(tx, q), lt)
	return lt.TriggerID, nil
}

func (tx *Tx) UpdateLabelTrigger(lt LabelTrigger) error {
	return tx.updateTrigger(lt.Trigger)
}
