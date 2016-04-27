package revere

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/yext/revere/probes"
)

type MonitorID int32

type Monitor struct {
	MonitorId   MonitorID
	Name        string
	Owner       string
	Description string
	Response    string
	ProbeType   probes.ProbeTypeId
	ProbeJson   string
	Changed     time.Time
	Version     int
	Archived    *time.Time
	Triggers    []*MonitorTrigger
	Labels      []*MonitorLabel
}

type MonitorTrigger struct {
	Trigger
	Subprobe string
	Delete   bool
}

type MonitorLabel struct {
	Label
	Subprobe string
	Create   bool
	Delete   bool
}

const (
	allMonitorFields        = "monitorid, name, owner, description, response, probetype, probe, changed, version, archived"
	allMonitorTriggerFields = "monitorid, subprobe, triggerid"
	allMonitorLabelFields   = "l.labelid, l.name, l.description, lm.subprobe"
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

func LoadMonitorsForLabel(db *sql.DB, labelId LabelID) ([]*Monitor, error) {
	rows, err := db.Query(
		`SELECT m.* FROM monitors m
		JOIN labels_monitors l on l.monitorid = m.monitorid
		WHERE l.labelid = ?
		ORDER BY name`,
		labelId)
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

func LoadMonitor(db *sql.DB, id MonitorID) (m *Monitor, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM monitors WHERE monitorid = %d", allMonitorFields, id))
	if err != nil {
		return nil, err
	}

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
	if err := rows.Scan(&m.MonitorId, &m.Name, &m.Owner, &m.Description, &m.Response, &m.ProbeType, &m.ProbeJson, &m.Changed, &m.Version, &m.Archived); err != nil {
		return nil, err
	}

	return &m, nil
}

func LoadMonitorLabels(db *sql.DB, monitorId MonitorID) ([]*MonitorLabel, error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM labels l
		JOIN labels_monitors lm on l.labelid=lm.labelid
		WHERE lm.monitorid = %d
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

func BatchLoadMonitorLabels(db *sql.DB, mIds []MonitorID) (map[MonitorID][]*MonitorLabel, error) {
	if len(mIds) == 0 {
		return nil, nil
	}

	monitors := make([]interface{}, len(mIds))
	for i, m := range mIds {
		monitors[i] = m
	}

	stmt, err := db.Prepare(fmt.Sprintf(`
		SELECT %s, lm.monitorid FROM labels l
		JOIN labels_monitors lm on l.labelid=lm.labelid
		WHERE lm.monitorid IN (%s)
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

	monitorLabels := make(map[MonitorID][]*MonitorLabel)
	for rows.Next() {
		var (
			ml  MonitorLabel
			mId MonitorID
		)
		if err := rows.Scan(&ml.LabelId, &ml.Name, &ml.Description, &ml.Subprobe, &mId); err != nil {
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
	if err := rows.Scan(&ml.LabelId, &ml.Name, &ml.Description, &ml.Subprobe); err != nil {
		return nil, err
	}

	return &ml, nil
}

func isExistingMonitor(db *sql.DB, id MonitorID) (exists bool) {
	if id == 0 {
		return false
	}

	err := db.QueryRow("SELECT EXISTS (SELECT * FROM monitors WHERE monitorid = ?)", id).Scan(&exists)
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
	if m.MonitorId == 0 {
		m.MonitorId, err = m.create(tx)
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
			err = t.save(tx, m.MonitorId)
			if err != nil {
				return
			}
		}
	}

	for _, ml := range m.Labels {
		err = ml.save(tx, m.MonitorId)
		if err != nil {
			return
		}
	}
	return
}

func (m *Monitor) create(tx *sql.Tx) (MonitorID, error) {
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
	return MonitorID(id), err
}

func (m *Monitor) update(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE monitors
		SET name=?, owner=?, description=?, response=?, probetype=?, probe=?, changed=?, version=version+1, archived=?
		WHERE monitorid=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	probeType, err := probes.ProbeTypeById(m.ProbeType)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(m.Name, m.Owner, m.Description, m.Response, probeType.Id(), m.ProbeJson, time.Now(), m.Archived, m.MonitorId)
	if err != nil {
		return err
	}
	return nil
}

func (mt *MonitorTrigger) Validate() (errs []string) {
	if err := validateSubprobe(mt.Subprobe); err != nil {
		errs = append(errs, err.Error())
	}
	errs = append(errs, mt.Trigger.Validate()...)
	return
}

func (mt *MonitorTrigger) save(tx *sql.Tx, mId MonitorID) (err error) {
	var newTriggerId TriggerID
	newTriggerId, err = mt.Trigger.save(tx)
	if err != nil {
		return
	}

	if mt.TriggerId == 0 {
		mt.TriggerId = newTriggerId
		err = mt.create(tx, mId)
	} else {
		err = mt.update(tx, mId)
	}
	return
}

func (mt *MonitorTrigger) create(tx *sql.Tx, mId MonitorID) error {
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO monitor_triggers(%s) VALUES (?, ?, ?)", allMonitorTriggerFields))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(mId, mt.Subprobe, mt.TriggerId)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (mt *MonitorTrigger) update(tx *sql.Tx, mId MonitorID) error {
	stmt, err := tx.Prepare(`
		UPDATE monitor_triggers
		SET subprobe=?
		WHERE monitorid=? AND triggerid=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(mt.Subprobe, mId, mt.TriggerId)
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
	if err := validateSubprobe(ml.Subprobe); err != nil {
		errs = append(errs, err.Error())
	}

	if !isExistingLabel(db, ml.LabelId) {
		errs = append(errs, fmt.Sprintf("Invalid label: %d", ml.LabelId))
	}
	return
}

func (ml *MonitorLabel) save(tx *sql.Tx, monitorId MonitorID) (err error) {
	if ml.Create {
		return ml.create(tx, monitorId)
	}
	if ml.Delete {
		return ml.delete(tx, monitorId)
	}
	return ml.update(tx, monitorId)
}

func (ml *MonitorLabel) create(tx *sql.Tx, monitorId MonitorID) error {
	stmt, err := tx.Prepare(
		`INSERT INTO labels_monitors(monitorid, labelid, subprobe) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(monitorId, ml.LabelId, ml.Subprobe)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (ml *MonitorLabel) update(tx *sql.Tx, monitorId MonitorID) error {
	stmt, err := tx.Prepare(`
		UPDATE labels_monitors
		SET subprobe=?
		WHERE monitorid=? AND labelid=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ml.Subprobe, monitorId, ml.LabelId)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (ml *MonitorLabel) delete(tx *sql.Tx, monitorId MonitorID) error {
	stmt, err := tx.Prepare(`
		DELETE FROM labels_monitors
		WHERE monitorid=? AND labelid=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(monitorId, ml.LabelId)
	if err != nil {
		return err
	}
	return stmt.Close()
}
