package db

import (
	"time"

	"github.com/juju/errors"

	"github.com/yext/revere/state"
)

type ReadingID int64

type Reading struct {
	ReadingID  ReadingID
	SubprobeID SubprobeID
	Recorded   time.Time
	State      state.State
}

func (db *DB) LoadReadings(subprobeID SubprobeID) ([]*Reading, error) {
	var readings []*Reading
	// TODO(eefi): Drop 2 suffix from table name.
	query := `SELECT * FROM pfx_readings2 WHERE subprobeid = ? ORDER BY recorded DESC`
	if err := db.Select(&readings, cq(db, query), subprobeID); err != nil {
		return nil, errors.Trace(err)
	}
	return readings, nil
}

func (tx *Tx) InsertReading(r Reading) error {
	// TODO(eefi): Drop 2 suffix from table name.
	q := `INSERT INTO pfx_readings2 (subprobeid, recorded, state)
	      VALUES (:subprobeid, :recorded, :state)`
	_, err := tx.NamedExec(cq(tx, q), r)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
