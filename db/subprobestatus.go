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

func (tx *Tx) InsertSubprobeStatus(s SubprobeStatus) error {
	q := `INSERT INTO pfx_subprobe_statuses (
	        subprobeid,
	        recorded,
	        state,
	        silenced,
	        enteredstate,
	        lastnormal
	      ) VALUES (
	        :subprobeid,
		:recorded,
		:state,
		:silenced,
		:enteredstate,
		:lastnormal
	      )`
	_, err := tx.NamedExec(cq(tx, q), s)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (tx *Tx) UpdateSubprobeStatus(s SubprobeStatus) error {
	q := `UPDATE pfx_subprobe_statuses
	      SET recorded = :recorded,
	          state = :state,
	          silenced = :silenced,
	          enteredstate = :enteredstate,
	          lastnormal = :lastnormal
	      WHERE subprobeid = :subprobeid`
	result, err := tx.NamedExec(cq(tx, q), s)
	if err != nil {
		return errors.Trace(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Trace(err)
	}

	if rows > 0 {
		return nil
	}

	var rowExists bool
	q = `SELECT EXISTS(SELECT 1 FROM pfx_subprobe_statuses WHERE subprobeid = ?)`
	err = tx.Get(&rowExists, cq(tx, q), s.SubprobeID)
	if err != nil {
		return errors.Trace(err)
	}

	if rowExists {
		// DB already had the data we were trying to write.
		return nil
	}

	// TODO(eefi): Try to fix the DB by inserting a new row?
	return errors.Errorf("no status row for subprobe with ID %d", s.SubprobeID)
}
