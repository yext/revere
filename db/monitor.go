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

func (db *DB) LoadAllMonitorIDs() ([]MonitorID, error) {
	return loadAllMonitorIDs(db)
}

func (tx *Tx) LoadAllMonitorIDs() ([]MonitorID, error) {
	return loadAllMonitorIDs(tx)
}

func loadAllMonitorIDs(dt dbOrTx) ([]MonitorID, error) {
	var ids []MonitorID
	err := dt.Select(&ids, cq(dt, "SELECT monitorid FROM pfx_monitors"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return ids, nil
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
