package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"
)

type MonitorID int32
type ProbeType int16

type Monitor struct {
	MonitorID   MonitorID
	Name        string
	Owner       string
	Description string
	Response    string
	ProbeType   ProbeType
	Probe       types.JSONText
	Changed     time.Time
	Version     int32
	Archived    *time.Time
}

type MonitorTrigger struct {
	MonitorID MonitorID
	// TODO(eefi): Rename column in DB to subprobes.
	Subprobes string `db:"subprobe"`
	*Trigger
}

type MonitorLabel struct {
	MonitorID MonitorID
	// TODO(eefi): Rename column in DB to subprobes.
	Subprobes string `db:"subprobe"`
	*Label
}

type MonitorVersionInfo struct {
	MonitorID MonitorID
	Version   int32
	Archived  *time.Time
}

func (tx *Tx) CreateMonitor(m *Monitor) (MonitorID, error) {
	q := `INSERT INTO pfx_monitors (name, owner, description, response, probetype, probe, changed, version, archived)
		VALUES (:name, :owner, :description, :response, :probetype, :probe, :changed, :version, :archived)`
	result, err := tx.NamedExec(cq(tx, q), m)
	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}
	return MonitorID(id), nil
}

func (tx *Tx) UpdateMonitor(m *Monitor) error {
	q := `UPDATE pfx_monitors
	      SET name=:name,
	          owner=:owner
	          description=:description
	          response=:response
	          probetype=:probetype
	          probe=:probe
	          changed=:changed
	          version=:version
	          archived=:archived
	      WHERE monitorid=:monitorid`
	_, err := tx.NamedExec(cq(tx, q), m)
	return errors.Trace(err)
}

func (db *DB) LoadMonitorVersionInfosUpdatedSince(t time.Time) ([]MonitorVersionInfo, error) {
	var infos []MonitorVersionInfo
	var err error

	q := "SELECT monitorid, version, archived FROM pfx_monitors"
	if t.IsZero() {
		err = db.Select(&infos, cq(db, q))
	} else {
		q += " WHERE changed >= ?"
		err = db.Select(&infos, cq(db, q), t)
	}
	if err != nil {
		return nil, errors.Trace(err)
	}

	return infos, nil
}

func (db *DB) LoadMonitor(id MonitorID) (*Monitor, error) {
	return loadMonitor(db, id)
}

func (tx *Tx) LoadMonitor(id MonitorID) (*Monitor, error) {
	return loadMonitor(tx, id)
}

func loadMonitor(dt dbOrTx, id MonitorID) (*Monitor, error) {
	dt = unsafe(dt)

	var m Monitor
	err := dt.Get(&m, cq(dt, "SELECT * FROM pfx_monitors WHERE monitorid = ?"), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &m, nil
}

func (db *DB) LoadMonitors() ([]*Monitor, error) {
	return loadMonitors(db)
}

func (tx *Tx) LoadMonitors() ([]*Monitor, error) {
	return loadMonitors(tx)
}

func loadMonitors(dt dbOrTx) ([]*Monitor, error) {
	var monitors []*Monitor
	if err := dt.Select(&monitors, cq(dt, "SELECT * FROM pfx_monitors ORDER BY name")); err != nil {
		return nil, errors.Trace(err)
	}
	return monitors, nil
}

func (tx *Tx) LoadMonitorsWithLabel(id LabelID) ([]*Monitor, error) {
	var monitors []*Monitor
	err := tx.Select(&monitors, cq(tx, `
		SELECT m.* FROM pfx_monitors m
		JOIN labels_monitors l on l.monitorid = m.monitorid
		WHERE l.labelid = ?
		ORDER BY name
	`), id)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return monitors, nil
}

func (db *DB) IsExistingMonitor(id MonitorID) (exists bool) {
	if id == 0 {
		return false
	}

	q := `SELECT EXISTS (SELECT * FROM pfx_monitors WHERE monitorid = ?)`
	err := db.Get(&exists, cq(db, q), id)
	if err != nil {
		return false
	}
	return
}

func (db *DB) LoadTriggersForMonitor(id MonitorID) ([]MonitorTrigger, error) {
	return loadTriggersForMonitor(db, id)
}

func (tx *Tx) LoadTriggersForMonitor(id MonitorID) ([]MonitorTrigger, error) {
	return loadTriggersForMonitor(tx, id)
}

func loadTriggersForMonitor(dt dbOrTx, id MonitorID) ([]MonitorTrigger, error) {
	dt = unsafe(dt)

	var mts []MonitorTrigger
	q := `SELECT *
	      FROM pfx_monitor_triggers
	      JOIN pfx_triggers USING (triggerid)
	      WHERE pfx_monitor_triggers.monitorid = ?`
	err := dt.Select(&mts, cq(dt, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return mts, nil
}

func (db *DB) LoadLabelsForMonitor(id MonitorID) ([]MonitorLabel, error) {
	return loadLabelsForMonitor(db, id)
}

func (tx *Tx) LoadLabelsForMonitor(id MonitorID) ([]MonitorLabel, error) {
	return loadLabelsForMonitor(tx, id)
}

func loadLabelsForMonitor(dt dbOrTx, id MonitorID) ([]MonitorLabel, error) {
	dt = unsafe(dt)

	var mls []MonitorLabel
	q := `SELECT *
	      FROM pfx_labels_monitors
	      JOIN pfx_labels USING (labelid)
	      WHERE pfx_labels_monitors.monitorid = ?`
	err := dt.Select(&mls, cq(dt, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return mls, nil
}

func (tx *Tx) BatchLoadMonitorLabels(mIDs []MonitorID) (map[MonitorID][]MonitorLabel, error) {
	if len(mIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`
		SELECT * FROM pfx_labels
		JOIN pfx_labels_monitors USING (labelid)
		WHERE monitorid IN (?)
	`, mIDs)
	if err != nil {
		return nil, err
	}
	// This isn't absolutely necessary but putting this here to support other potential backends
	query = tx.Rebind(query)

	rows, err := tx.Queryx(cq(tx, query), args...)
	if err != nil {
		return nil, err
	}

	monitorLabels := make(map[MonitorID][]MonitorLabel)
	for rows.Next() {
		var ml MonitorLabel
		if err = rows.StructScan(&ml); err != nil {
			return nil, err
		}
		monitorLabels[ml.MonitorID] = append(monitorLabels[ml.MonitorID], ml)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return monitorLabels, nil
}

func (tx *Tx) CreateMonitorTrigger(mt MonitorTrigger) (TriggerID, error) {
	var err error
	mt.TriggerID, err = tx.createTrigger(mt.Trigger)
	if err != nil {
		return 0, errors.Trace(err)
	}

	// TODO(psingh): Change field to subprobe once done renaming field
	q := `INSERT INTO pfx_monitor_triggers (monitorid, subprobe, triggerid)
	      VALUES (:monitorid, :subprobe, :triggerid)`
	_, err = tx.NamedExec(cq(tx, q), mt)
	return mt.TriggerID, errors.Trace(err)
}

func (tx *Tx) UpdateMonitorTrigger(mt MonitorTrigger) error {
	err := tx.updateTrigger(mt.Trigger)
	if err != nil {
		return errors.Trace(err)
	}

	// TODO(psingh): Change field to subprobe once done renaming field
	q := `UPDATE pfx_monitor_triggers
	      SET subprobe=:subprobe
	      WHERE triggerid=:triggerid`
	_, err = tx.NamedExec(cq(tx, q), mt)
	return errors.Trace(err)
}

func (tx *Tx) DeleteMonitorTrigger(triggerID TriggerID) error {
	return tx.deleteTrigger(triggerID)
}

func (tx *Tx) CreateMonitorLabel(ml MonitorLabel) error {
	// TODO(psingh): Change field to subprobe once done renaming field
	q := `INSERT INTO pfx_labels_monitors (labelid, monitorid, subprobe)
	      VALUES (:labelid, :monitorid, :subprobe)`
	_, err := tx.NamedExec(cq(tx, q), ml)
	return errors.Trace(err)
}

func (tx *Tx) UpdateMonitorLabel(ml MonitorLabel) error {
	// TODO(psingh): Change field to subprobe once done renaming field
	q := `UPDATE pfx_labels_monitors
	      SET subprobe=:subprobe
	      WHERE labelid=:labelid AND monitorid=:monitorid`
	_, err := tx.NamedExec(cq(tx, q), ml)
	return errors.Trace(err)
}

func (tx *Tx) DeleteMonitorLabel(ml MonitorLabel) error {
	// TODO(psingh): Change field to subprobe once done renaming field
	q := `DELETE FROM pfx_labels_monitors
	      WHERE labelid=:labelid AND monitorid=:monitorid`
	_, err := tx.NamedExec(cq(tx, q), ml)
	return errors.Trace(err)
}
