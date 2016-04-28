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
	q := `INSERT INTO pfx_readings (subprobeid, recorded, state)
	      VALUES (:subprobeid, :recorded, :state)`
	_, err := tx.NamedExec(cq(tx, q), r)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
