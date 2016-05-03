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
