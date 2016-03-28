package revere

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yext/revere/probes"
)

type Monitor struct {
	Id          uint               `json:"id,omitempty"`
	Name        string             `json:"name"`
	Owner       string             `json:"owner"`
	Description string             `json:"description"`
	Response    string             `json:"response"`
	ProbeType   probes.ProbeTypeId `json:"probeType"`
	ProbeJson   string             `json:"probe"`
	Changed     time.Time          `json:"-"`
	Version     int                `json:"-"`
	Archived    *time.Time         `json:"-"` // nullable
	Triggers    []*MonitorTrigger  `json:"triggers,omitempty"`
}

type MonitorTrigger struct {
	Trigger
	Subprobes string `json:"subprobes"`
	Delete    bool   `json:"delete,omitempty"`
}

const (
	allMonitorFields        = "id, name, owner, description, response, probeType, probe, changed, version, archived"
	allMonitorTriggerFields = "monitor_id, subprobe, trigger_id"
)

func (m *Monitor) Validate() (errs []string) {
	if m.Name == "" {
		errs = append(errs, fmt.Sprintf("Monitor name is required"))
	}

	probeType, err := probes.ProbeTypeById(m.ProbeType)
	if err != nil {
		errs = append(errs, err.Error())
	}
	probe, err := probeType.Load(m.ProbeJson)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Invalid probe for monitor: %s", m.ProbeJson))
	}
	errs = append(errs, probe.Validate()...)

	for _, mt := range m.Triggers {
		errs = append(errs, mt.Validate()...)
	}
	return
}

func LoadMonitors(db *sql.DB) ([]*Monitor, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM monitors ORDER BY name", allMonitorFields))
	if err != nil {
		return nil, err
	}

	return loadMonitorsFromRows(rows)
}

func LoadMonitorsForLabel(db *sql.DB, labelId uint) ([]*Monitor, error) {
	rows, err := db.Query(fmt.Sprintf(
		`SELECT %s FROM monitors
		JOIN labels_monitors lt on lt.monitor_id = id
		WHERE lt.label_id = %d ORDER BY name`,
		allMonitorFields, labelId))
	if err != nil {
		return nil, err
	}

	return loadMonitorsFromRows(rows)
}

func loadMonitorsFromRows(rows *sql.Rows) ([]*Monitor, error) {
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

		m.Triggers, err = LoadMonitorTriggers(db, id)
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
	if err := rows.Scan(&m.Id, &m.Name, &m.Owner, &m.Description, &m.Response, &m.ProbeType, &m.ProbeJson, &m.Changed, &m.Version, &m.Archived); err != nil {
		return nil, err
	}

	return &m, nil
}

func isExistingMonitor(db *sql.DB, id uint) (exists bool) {
	if id == 0 {
		return false
	}

	err := db.QueryRow("SELECT EXISTS (SELECT * FROM monitors WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return
}

func (m *Monitor) Save(db *sql.DB) (err error) {
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
	if m.Id == 0 {
		m.Id, err = m.create(tx)
	} else {
		err = m.update(tx)
	}
	if err != nil {
		return
	}

	// Create/Update/Delete Monitor Triggers
	for _, t := range m.Triggers {
		if t.Delete {
			err = t.delete(tx)
			if err != nil {
				return
			}
		} else {
			err = t.save(tx, m.Id)
			if err != nil {
				return
			}
		}
	}
	return
}

func (m *Monitor) create(tx *sql.Tx) (uint, error) {
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO monitors(%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", allMonitorFields))
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	probeType, err := probes.ProbeTypeById(m.ProbeType)
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

func (m *Monitor) update(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE monitors
		SET name=?, owner=?, description=?, response=?, probeType=?, probe=?, changed=?, version=version+1, archived=?
		WHERE id=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	probeType, err := probes.ProbeTypeById(m.ProbeType)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(m.Name, m.Owner, m.Description, m.Response, probeType.Id(), m.ProbeJson, time.Now(), m.Archived, m.Id)
	if err != nil {
		return err
	}
	return nil
}

func (mt *MonitorTrigger) Validate() (errs []string) {
	if err := validateSubprobes(mt.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}
	errs = append(errs, mt.Trigger.Validate()...)
	return
}

func (mt *MonitorTrigger) save(tx *sql.Tx, mId uint) (err error) {
	var newTriggerId uint
	newTriggerId, err = mt.Trigger.save(tx)
	if err != nil {
		return
	}

	if mt.Id == 0 {
		mt.Id = newTriggerId
		err = mt.create(tx, mId)
	} else {
		err = mt.update(tx, mId)
	}
	return
}

func (mt *MonitorTrigger) create(tx *sql.Tx, mId uint) error {
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO monitor_triggers(%s) VALUES (?, ?, ?)", allMonitorTriggerFields))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(mId, mt.Subprobes, mt.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (mt *MonitorTrigger) update(tx *sql.Tx, mId uint) error {
	stmt, err := tx.Prepare(`
		UPDATE monitor_triggers
		SET subprobe=?
		WHERE monitor_id=? AND trigger_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(mt.Subprobes, mId, mt.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (mt *MonitorTrigger) delete(tx *sql.Tx) error {
	// Trigger delete will cascade to monitor triggers
	return mt.Trigger.delete(tx)
}
