package revere

import (
	"database/sql"
	"fmt"
	"time"
)

const allReadingFields = "r.id, r.subprobe_id, r.state, r.recorded"

type Reading struct {
	Id         uint
	SubprobeId uint
	State      State
	StateStr   string
	Recorded   time.Time
}

func LoadReadings(db *sql.DB, subprobeId uint) (readings []*Reading, err error) {
	rows, err := db.Query(
		fmt.Sprintf("SELECT %s FROM readings2 r WHERE r.subprobe_id = %d ORDER BY r.recorded DESC", allReadingFields, subprobeId))
	if err != nil {
		return nil, err
	}
	readings = make([]*Reading, 0)
	for rows.Next() {
		s, err := loadReadingFromRow(rows)
		if err != nil {
			return nil, err
		}
		readings = append(readings, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return readings, nil
}

func loadReadingFromRow(rows *sql.Rows) (*Reading, error) {
	var r Reading
	if err := rows.Scan(&r.Id, &r.SubprobeId, &r.State, &r.Recorded); err != nil {
		return nil, err
	}

	r.StateStr = States(r.State)

	return &r, nil
}
