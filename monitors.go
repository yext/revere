package revere

import (
	"database/sql"
	"time"
)

type Monitor struct {
	Id          uint
	Name        string
	Owner       string
	Description string
	Response    string
	ProbeType   int
	Probe       string
	Changed     time.Time
	Version     int
	Archived    *time.Time // nullable
}

func LoadMonitors(db *sql.DB) (map[uint]Monitor, error) {
	rows, err := db.Query("SELECT id, name, owner, description, response, probeType, probe, changed, version, archived FROM monitors")
	if err != nil {
		return nil, err
	}

	allMonitors := make(map[uint]Monitor)
	for rows.Next() {
		m, err := LoadMonitor(rows)
		if err != nil {
			return nil, err
		}
		allMonitors[m.Id] = *m
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allMonitors, nil
}

func LoadMonitor(rows *sql.Rows) (*Monitor, error) {
	var m Monitor
	if err := rows.Scan(&m.Id, &m.Name, &m.Owner, &m.Description, &m.Response, &m.ProbeType, &m.Probe, &m.Changed, &m.Version, &m.Archived); err != nil {
		return nil, err
	}
	return &m, nil
}
