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
	var silences []*MonitorSilence
	q := `SELECT s.*, m.name AS monitorname
		  FROM pfx_silences s
		  JOIN pfx_monitors m USING (monitorid)`
	if err := db.Select(&silences, cq(db, q)); err != nil {
		return nil, errors.Trace(err)
	}
	return silences, nil
}
