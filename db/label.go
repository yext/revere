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
	dt = unsafe(dt)

	labels := []*Label{}
	err := dt.Select(&labels, cq(dt, "SELECT * FROM pfx_labels ORDER BY name"))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return labels, nil
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
	      WHERE pfx_labels_monitors.labelid = ?`
	err := dt.Select(&lms, cq(dt, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return lms, nil
}

type LabelTriggerWithSubprobes struct {
	LabelTrigger
	// TODO(eefi): Rename column in DB to subprobes.
	Subprobes string `db:"subprobe"`
}

func (tx *Tx) LoadLabelTriggersForMonitor(id MonitorID) ([]LabelTriggerWithSubprobes, error) {
	var results []LabelTriggerWithSubprobes
	// TODO(eefi): Update "subprobe" when DB column is renamed.
	q := `SELECT labelid, triggerid, pfx_triggers.*, pfx_labels_monitors.subprobe
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

func (tx *Tx) CreateLabelTrigger(lt LabelTrigger) error {
	var err error
	lt.TriggerID, err = tx.createTrigger(lt.Trigger)
	if err != nil {
		return errors.Trace(err)
	}

	// TODO(psingh): Change field to subprobe once done renaming field
	q := `INSERT INTO pfx_label_triggers (labelid, triggerid)
	      VALUES (:labelid, :triggerid)`
	_, err = tx.NamedExec(cq(tx, q), lt)
	return errors.Trace(err)
}

func (tx *Tx) UpdateLabelTrigger(lt LabelTrigger) error {
	return tx.updateTrigger(lt.Trigger)
}
