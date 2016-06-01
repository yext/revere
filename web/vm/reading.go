package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
)

type Reading struct {
	*revere.Reading
}

func (r *Reading) Id() int64 {
	return int64(r.ReadingId)
}

func AllReadingsFromSubprobe(db *sql.DB, id revere.SubprobeID) ([]*Reading, error) {
	rs, err := revere.LoadReadings(db, id)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, fmt.Errorf("No readings found for subprobe: %d", id)
	}

	return newReadingsFromModel(db, rs), nil
}

func newReadingFromModel(db *sql.DB, r *revere.Reading) *Reading {
	viewmodel := new(Reading)
	viewmodel.Reading = r

	return viewmodel
}

func newReadingsFromModel(db *sql.DB, rs []*revere.Reading) []*Reading {
	readings := make([]*Reading, len(rs))
	for i, r := range rs {
		readings[i] = newReadingFromModel(db, r)
	}
	return readings
}

func BlankReading(db *sql.DB) *Reading {
	viewmodel := new(Reading)
	viewmodel.Reading = new(revere.Reading)

	return viewmodel
}
