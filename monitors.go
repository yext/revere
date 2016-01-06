package revere

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/yext/revere/probes"
)

type Monitor struct {
	Id            uint               `json:"id,omitempty"`
	Name          string             `json:"name"`
	Owner         string             `json:"owner"`
	Description   string             `json:"description"`
	Response      string             `json:"response"`
	Probe         probes.Probe       `json:"-"`
	ProbeType     probes.ProbeTypeId `json:"probeType"`
	ProbeJson     string             `json:"probe"`
	ProbeTemplate template.HTML      `json:"-"`
	Changed       time.Time          `json:"-"`
	Version       int                `json:"-"`
	Archived      *time.Time         `json:"-"` // nullable
	Triggers      []*Trigger         `json:"triggers"`
}

const allMonitorFields = "id, name, owner, description, response, probeType, probe, changed, version, archived"

func (m *Monitor) Validate() (errs []string) {
	if m.Name == "" {
		errs = append(errs, fmt.Sprintf("Monitor name is required"))
	}

	probeType, err := probes.GetProbeType(m.ProbeType)
	if err != nil {
		errs = append(errs, err.Error())
	}
	probe, err := probeType.Load(m.ProbeJson)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Invalid probe for monitor: %s", m.ProbeJson))
	}
	errs = append(errs, probe.Validate()...)

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

	// Load Triggers
	m.Triggers, err = LoadTriggers(db, id)
	if err != nil {
		return nil, err
	}

	// Load Probe
	probeType, err := probes.GetProbeType(m.ProbeType)
	if err != nil {
		return nil, err
	}

	m.Probe, err = probeType.Load(m.ProbeJson)
	if err != nil {
		return nil, err
	}

	m.ProbeTemplate, err = m.Probe.Render()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func loadMonitorFromRow(rows *sql.Rows) (*Monitor, error) {
	var m Monitor
	if err := rows.Scan(&m.Id, &m.Name, &m.Owner, &m.Description, &m.Response, &m.ProbeType, &m.ProbeJson, &m.Changed, &m.Version, &m.Archived); err != nil {
		return nil, err
	}

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

	probeType, err := probes.GetProbeType(m.ProbeType)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(nil, m.Name, m.Owner, m.Description, m.Response, probeType.Id(), m.ProbeJson, time.Now(), 1, m.Archived)
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

	probeType, err := probes.GetProbeType(m.ProbeType)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(m.Name, m.Owner, m.Description, m.Response, probeType.Id(), m.ProbeJson, time.Now(), m.Archived, m.Id)
	if err != nil {
		return err
	}
	return nil
}
