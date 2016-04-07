package revere

import (
	"database/sql"
	"fmt"
	"strings"
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
	Labels      []*MonitorLabel    `json:"labels,omitempty"`
}

type MonitorTrigger struct {
	Trigger
	Subprobes string `json:"subprobes"`
	Delete    bool   `json:"delete,omitempty"`
}

type MonitorLabel struct {
	Label
	Subprobes string `json:"subprobes"`
	Create    bool   `json:"create,omitempty"`
	Delete    bool   `json:"delete,omitempty"`
}

const (
	allMonitorFields        = "id, name, owner, description, response, probeType, probe, changed, version, archived"
	allMonitorTriggerFields = "monitor_id, subprobe, trigger_id"
	allMonitorLabelFields   = "l.Id, l.Name, l.Description, lm.Subprobes"
)

func (m *Monitor) Validate(db *sql.DB) (errs []string) {
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

	for _, ml := range m.Labels {
		errs = append(errs, ml.Validate(db)...)
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

		m.Labels, err = LoadMonitorLabels(db, id)
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

func LoadMonitorLabels(db *sql.DB, monitorId uint) ([]*MonitorLabel, error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM labels l
		JOIN labels_monitors lm on l.id=lm.label_id
		WHERE lm.monitor_id = %d
	`, allMonitorLabelFields, monitorId))
	if err != nil {
		return nil, err
	}

	monitorLabels := make([]*MonitorLabel, 0)
	for rows.Next() {
		ml, err := loadMonitorLabelsFromRow(rows)
		if err != nil {
			return nil, err
		}
		monitorLabels = append(monitorLabels, ml)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return monitorLabels, nil
}

func BatchLoadMonitorLabels(db *sql.DB, mIds []uint) (map[uint][]*MonitorLabel, error) {
	monitors := make([]interface{}, len(mIds))
	for i, m := range mIds {
		monitors[i] = m
	}

	stmt, err := db.Prepare(fmt.Sprintf(`
		SELECT %s, lm.monitor_id FROM labels l
		JOIN labels_monitors lm on l.id=lm.label_id
		WHERE lm.monitor_id IN (%s)
	`, allMonitorLabelFields,
		strings.TrimSuffix(strings.Repeat("?,", len(monitors)), ",")))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(monitors...)
	if err != nil {
		return nil, err
	}

	monitorLabels := make(map[uint][]*MonitorLabel)
	for rows.Next() {
		var (
			ml  MonitorLabel
			mId uint
		)
		if err := rows.Scan(&ml.Id, &ml.Name, &ml.Description, &ml.Subprobes, &mId); err != nil {
			return nil, err
		}
		monitorLabels[mId] = append(monitorLabels[mId], &ml)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return monitorLabels, nil
}

func loadMonitorLabelsFromRow(rows *sql.Rows) (*MonitorLabel, error) {
	var ml MonitorLabel
	if err := rows.Scan(&ml.Id, &ml.Name, &ml.Description, &ml.Subprobes); err != nil {
		return nil, err
	}

	return &ml, nil
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

	for _, ml := range m.Labels {
		err = ml.save(tx, m.Id)
		if err != nil {
			return
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

func (ml *MonitorLabel) Validate(db *sql.DB) (errs []string) {
	if err := validateSubprobes(ml.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}

	if !isExistingLabel(db, ml.Id) {
		errs = append(errs, fmt.Sprintf("Invalid label: %d", ml.Id))
	}
	return
}

func (ml *MonitorLabel) save(tx *sql.Tx, monitorId uint) (err error) {
	if ml.Create {
		return ml.create(tx, monitorId)
	}
	if ml.Delete {
		return ml.delete(tx, monitorId)
	}
	return ml.update(tx, monitorId)
}

func (ml *MonitorLabel) create(tx *sql.Tx, monitorId uint) error {
	stmt, err := tx.Prepare(
		`INSERT INTO labels_monitors(monitor_id, label_id, subprobes) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(monitorId, ml.Id, ml.Subprobes)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (ml *MonitorLabel) update(tx *sql.Tx, monitorId uint) error {
	stmt, err := tx.Prepare(`
		UPDATE labels_monitors
		SET subprobes=?
		WHERE monitor_id=? AND label_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ml.Subprobes, monitorId, ml.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (ml *MonitorLabel) delete(tx *sql.Tx, monitorId uint) error {
	stmt, err := tx.Prepare(`
		DELETE FROM labels_monitors
		WHERE monitor_id=? AND label_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(monitorId, ml.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}
