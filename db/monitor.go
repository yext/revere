package db

import (
	"database/sql"
	"time"

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
	Subprobes string
	*Label
}

type MonitorVersionInfo struct {
	MonitorID MonitorID
	Version   int32
	Archived  *time.Time
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
	dt = unsafe(dt)

	monitors := []*Monitor{}
	err := dt.Select(&monitors, cq(dt, "SELECT * FROM pfx_monitors ORDER BY name"))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return monitors, nil
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

func (tx *Tx) CreateMonitorTrigger(mt MonitorTrigger) error {
	var err error
	mt.TriggerID, err = tx.createTrigger(mt.Trigger)
	if err != nil {
		return errors.Trace(err)
	}

	// TODO(psingh): Change field to subprobe once done renaming field
	q := `INSERT INTO pfx_monitor_triggers (monitorid, subprobe, triggerid)
	      VALUES (:monitorid, :subprobes, :triggerid)`
	_, err = tx.NamedExec(cq(tx, q), mt)
	return errors.Trace(err)
}

func (tx *Tx) UpdateMonitorTrigger(mt MonitorTrigger) error {
	err := tx.updateTrigger(mt.Trigger)
	if err != nil {
		return errors.Trace(err)
	}

	// TODO(psingh): Change field to subprobe once done renaming field
	q := `UPDATE pfx_monitor_triggers
	      SET subprobe=:subprobes
	      WHERE triggerid=:triggerid`
	_, err = tx.NamedExec(cq(tx, q), mt)
	return errors.Trace(err)
}
