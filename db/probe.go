package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"
)

type ProbeInfo struct {
	Probe     types.JSONText
	ProbeType ProbeType
}

func (db *DB) LoadProbeByMonitorID(id MonitorID) (*ProbeInfo, error) {
	return loadProbeByMonitorID(db, id)
}

func (tx *Tx) LoadProbeByMonitorID(id MonitorID) (*ProbeInfo, error) {
	return loadProbeByMonitorID(tx, id)
}

func loadProbeByMonitorID(dt dbOrTx, id MonitorID) (*ProbeInfo, error) {
	dt = unsafe(dt)

	var p ProbeInfo
	err := dt.Get(&p, cq(dt, "SELECT probe, probetype FROM pfx_monitors WHERE monitorid = ?"), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &p, nil
}
