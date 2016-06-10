package db

import (
	"database/sql"
	"time"

	"github.com/juju/errors"
)

type SilenceID int32

type Silence struct {
	SilenceID SilenceID
	MonitorID MonitorID
	// TODO(eefi): Rename column in DB to subprobes.
	Subprobes string `db:"subprobe"`
	Start     time.Time
	End       time.Time
}

type MonitorSilence struct {
	MonitorName string
	*Silence
}

func (db *DB) IsExistingSilence(id SilenceID) (exists bool) {
	if id == 0 {
		return false
	}

	q := `SELECT EXISTS (SELECT * FROM pfx_silences WHERE silenceid = ?)`
	err := db.Get(&exists, cq(db, q), id)
	if err != nil {
		return false
	}
	return
}

func (tx *Tx) CreateMonitorSilence(monitorSilence *MonitorSilence) (SilenceID, error) {
	q := `INSERT INTO pfx_silences (monitorid, subprobe, start, end)
	VALUES (:monitorid, :subprobe, :start, :end)`
	result, err := tx.NamedExec(cq(tx, q), monitorSilence)
	if err != nil {
		return 0, errors.Trace(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}
	return SilenceID(id), nil
}

func (tx *Tx) UpdateMonitorSilence(monitorSilence *MonitorSilence) error {
	q := `UPDATE pfx_silences
	     SET monitorid=:monitorid, subprobe=:subprobe, start=:start, end=:end
		 WHERE silenceid=:silenceid`
	_, err := tx.NamedExec(cq(tx, q), monitorSilence)
	return errors.Trace(err)
}

func (db *DB) LoadActiveSilencesForMonitor(monitorID MonitorID) ([]Silence, error) {
	var silences []Silence
	q := `SELECT * FROM pfx_silences
	      WHERE monitorid = ? AND start <= UTC_TIMESTAMP() AND UTC_TIMESTAMP() <= end`
	if err := db.Select(&silences, cq(db, q), monitorID); err != nil {
		return nil, errors.Trace(err)
	}
	return silences, nil
}

func (db *DB) LoadMonitorSilence(id SilenceID) (*MonitorSilence, error) {
	return loadMonitorSilence(db, id)
}

func (tx *Tx) LoadMonitorSilence(id SilenceID) (*MonitorSilence, error) {
	return loadMonitorSilence(tx, id)
}

func loadMonitorSilence(dt dbOrTx, id SilenceID) (*MonitorSilence, error) {
	var s MonitorSilence
	q := `SELECT s.*, m.name AS monitorname
		  FROM pfx_silences s
		  JOIN pfx_monitors m USING (monitorid)
		  WHERE s.silenceid = ?`
	if err := dt.Get(&s, cq(dt, q), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &s, nil
}

func (db *DB) LoadMonitorSilences() ([]*MonitorSilence, error) {
	return loadMonitorSilences(db)
}

func (tx *Tx) LoadMonitorSilences() ([]*MonitorSilence, error) {
	return loadMonitorSilences(tx)
}

func loadMonitorSilences(dt dbOrTx) ([]*MonitorSilence, error) {
	//TODO(fchen): maybe put LIMIT or only filter for active silences because this could return quite a few
	var silences []*MonitorSilence
	q := `SELECT s.*, m.name AS monitorname
		  FROM pfx_silences s
		  JOIN pfx_monitors m USING (monitorid)`
	if err := dt.Select(&silences, cq(dt, q)); err != nil {
		return nil, errors.Trace(err)
	}
	return silences, nil
}
