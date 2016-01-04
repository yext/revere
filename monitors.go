package revere

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Monitor struct {
	Id          uint       `json:"id,omitempty"`
	Name        string     `json:"name"`
	Owner       string     `json:"owner"`
	Description string     `json:"description"`
	Response    string     `json:"response"`
	ProbeType   string     `json:"probeType"`
	Probe       string     `json:"probe"`
	Changed     time.Time  `json:"-"`
	Version     int        `json:"-"`
	Archived    *time.Time `json:"-"` // nullable
	Triggers    []*Trigger `json:"triggers"`
}

const allMonitorFields = "id, name, owner, description, response, probeType, probe, changed, version, archived"

// Probe types
type ProbeType int

const (
	graphiteThreshold ProbeType = iota
)

var ProbeTypes = map[ProbeType]string{
	graphiteThreshold: "Graphite Threshold",
}

var reverseProbeTypes map[string]ProbeType

func init() {
	reverseProbeTypes = make(map[string]ProbeType)
	for k, v := range ProbeTypes {
		reverseProbeTypes[v] = k
	}
}

func (m *Monitor) Validate() (errs []string) {
	if m.Name == "" {
		errs = append(errs, fmt.Sprintf("Monitor name is required"))
	}

	if _, ok := reverseProbeTypes[m.ProbeType]; !ok {
		errs = append(errs, fmt.Sprintf("Invalid probe type for monitor: %s", m.ProbeType))
	}

	// TODO(psingh): Add proper validation once we implement the ui for probes
	var probeJson interface{}
	if err := json.Unmarshal([]byte(m.Probe), &probeJson); err != nil {
		errs = append(errs, fmt.Sprintf("Invalid json for probe"))
	}

	for _, t := range m.Triggers {
		errs = append(errs, t.Validate()...)
	}
	return
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

	return &m, nil
}

func (m *Monitor) SaveMonitor(db *sql.DB) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Create/Update Monitor
	// TODO: change to int64
	if m.Id == 0 {
		m.Id, err = m.createMonitor(tx)
	} else {
		err = m.updateMonitor(tx)
	}
	if err != nil {
		return
	}

	// Create/Update Triggers
	for _, t := range m.Triggers {
		err = t.saveTrigger(tx, m)
		if err != nil {
			return
		}
	}
	return
}

func (m *Monitor) createMonitor(tx *sql.Tx) (uint, error) {
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO monitors(%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", allMonitorFields))
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(nil, m.Name, m.Owner, m.Description, m.Response, reverseProbeTypes[m.ProbeType], m.Probe, time.Now().UTC(), 1, m.Archived)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint(id), err
}

func (m *Monitor) updateMonitor(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE monitors
		SET name=?, owner=?, description=?, response=?, probeType=?, probe=?, changed=?, version=version+1, archived=?
		WHERE id=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Name, m.Owner, m.Description, m.Response, reverseProbeTypes[m.ProbeType], m.Probe, time.Now().UTC(), m.Archived, m.Id)
	if err != nil {
		return err
	}
	return nil
}
