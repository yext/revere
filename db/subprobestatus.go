package db

import (
	"time"

	"github.com/juju/errors"

	"github.com/yext/revere/state"
)

type SubprobeStatus struct {
	SubprobeID   SubprobeID
	Recorded     time.Time
	State        state.State
	Silenced     bool
	EnteredState time.Time
	LastNormal   time.Time
}

func (db *DB) LoadSubprobeStatusesForMonitor(id MonitorID) (map[string]SubprobeStatus, error) {
	return loadSubprobeStatusesForMonitor(db, id)
}

func (tx *Tx) LoadSubprobeStatusesForMonitor(id MonitorID) (map[string]SubprobeStatus, error) {
	return loadSubprobeStatusesForMonitor(tx, id)
}

func loadSubprobeStatusesForMonitor(dt dbOrTx, id MonitorID) (map[string]SubprobeStatus, error) {
	dt = unsafe(dt)

	var data []struct {
		Name string
		SubprobeStatus
	}
	q := `SELECT pfx_subprobes.name, pfx_subprobe_statuses.*
	      FROM pfx_subprobes
	      JOIN pfx_subprobe_statuses USING (subprobeid)
	      WHERE pfx_subprobes.monitorid = ?`
	err := dt.Select(&data, cq(dt, q), id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result := make(map[string]SubprobeStatus)
	for _, d := range data {
		result[d.Name] = d.SubprobeStatus
	}

	return result, nil
}
