package db

import (
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

func (db *DB) LoadActiveSilencesForMonitor(monitorID MonitorID) ([]Silence, error) {
	var silences []Silence
	q := `SELECT * FROM pfx_silences
	      WHERE monitorid = ? AND start <= UTC_TIMESTAMP() AND UTC_TIMESTAMP() <= end`
	if err := db.Select(&silences, cq(db, q), monitorID); err != nil {
		return nil, errors.Trace(err)
	}
	return silences, nil
}
