package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type Monitor struct {
	Id          uint
	Name        string
	Owner       string
	Description string
	Response    string
	ProbeType   string
	Probe       string
	Changed     time.Time
	Version     int
	Archived    *time.Time // nullable
}

const allMonitorFields = "id, name, owner, description, response, probeType, probe, changed, version, archived"

// Probe types
type ProbeType int

const (
	graphiteThreshold = iota
)

var ProbeTypes = map[ProbeType]string{
	graphiteThreshold: "Graphite Threshold",
}

func LoadMonitors(db *sql.DB) ([]*Monitor, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM monitors ORDER BY name", allMonitorFields))
	if err != nil {
		return nil, err
	}

	allMonitors := make([]*Monitor, 0)
	for rows.Next() {
		m, err := loadMonitorFromRow(rows)
		if err != nil {
			return nil, err
		}
		allMonitors = append(allMonitors, m)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allMonitors, nil
}

func LoadMonitor(db *sql.DB, id uint) (m *Monitor, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM monitors WHERE id = %d", allMonitorFields, id))
	if rows.Next() {
		m, err = loadMonitorFromRow(rows)
		if err != nil {
			return nil, err
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func loadMonitorFromRow(rows *sql.Rows) (*Monitor, error) {
	var m Monitor
	var pt ProbeType
	if err := rows.Scan(&m.Id, &m.Name, &m.Owner, &m.Description, &m.Response, &pt, &m.Probe, &m.Changed, &m.Version, &m.Archived); err != nil {
		return nil, err
	}

	m.ProbeType = ProbeTypes[pt]
	m.Changed = ChangeLoc(m.Changed, time.UTC)
	if m.Archived != nil {
		t := ChangeLoc(*m.Archived, time.UTC)
		m.Archived = &t
	}

	return &m, nil
}

func ChangeLoc(t time.Time, l *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), l)
}
