package vm

import (
	"database/sql"
	"time"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
)

type Reading struct {
	ReadingID  db.ReadingID
	SubprobeID db.SubprobeID
	State      state.State
	StateStr   string
	Recorded   time.Time
}

func (r *Reading) Id() int64 {
	return int64(r.ReadingID)
}

func AllReadingsFromSubprobe(db *sql.DB, id db.SubprobeID) ([]*Reading, error) {
	rs, err := db.LoadReadings(db, id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if rs == nil {
		return nil, errors.Errorf("No readings found for subprobe: %d", id)
	}

	return newReadingsFromModel(rs), nil
}

func newReadingFromModel(reading *db.Reading) *Reading {
	return &Reading{
		ReadingID:  reading.ReadingID,
		SubprobeID: reading.SubprobeID,
		State:      reading.State,
		StateStr:   reading.State.String(),
		Recorded:   reading.Recorded,
	}
}

func newReadingsFromModel(readings []*db.Reading) []*Reading {
	rs := make([]*Reading, len(readings))
	for i, reading := range readings {
		rs[i] = newReadingFromModel(reading)
	}
	return rs
}

func BlankReading() *Reading {
	return &Reading{}
}
