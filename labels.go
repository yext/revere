package revere

import (
	"database/sql"
	"fmt"
)

type LabelID int32

type Label struct {
	LabelId     LabelID         `json:"id,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Triggers    []*LabelTrigger `json:"triggers,omitempty"`
	Monitors    []*LabelMonitor `json:"monitors,omitempty"`
}

type LabelMonitor struct {
	Monitor
	Subprobe string `json:"subprobe"`
	Create   bool   `json:"create,omitempty"`
	Delete   bool   `json:"delete,omitempty"`
}

type LabelTrigger struct {
	Trigger
	Delete bool `json:"delete,omitempty"`
}

const (
	allLabelFields        = "labelid, name, description"
	allLabelMonitorFields = "m.monitorid, m.name, m.description, lm.subprobe"
	allLabelTriggerFields = "labelid, triggerid"
)

func LoadLabels(db *sql.DB) (labels []*Label, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM labels ORDER BY name", allLabelFields))
	if err != nil {
		return nil, err
	}
	allLabels := make([]*Label, 0)
	for rows.Next() {
		m, err := loadLabelFromRow(rows)
		if err != nil {
			return nil, err
		}
		allLabels = append(allLabels, m)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allLabels, nil
}

func LoadLabel(db *sql.DB, id LabelID) (l *Label, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM labels WHERE labelid = %d", allLabelFields, id))
	if rows.Next() {
		l, err = loadLabelFromRow(rows)
		if err != nil {
			return nil, err
		}

		l.Triggers, err = LoadLabelTriggers(db, id)
		if err != nil {
			return nil, err
		}

		l.Monitors, err = LoadLabelMonitors(db, id)
		if err != nil {
			return nil, err
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return l, nil
}

func loadLabelFromRow(rows *sql.Rows) (*Label, error) {
	var l Label
	if err := rows.Scan(&l.LabelId, &l.Name, &l.Description); err != nil {
		return nil, err
	}

	return &l, nil
}

func LoadLabelMonitors(db *sql.DB, labelId LabelID) ([]*LabelMonitor, error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM monitors m
		JOIN labels_monitors lm on m.monitorid=lm.monitorid
		WHERE lm.labelid = %d
	`, allLabelMonitorFields, labelId))
	if err != nil {
		return nil, err
	}

	labelMonitors := make([]*LabelMonitor, 0)
	for rows.Next() {
		var lm LabelMonitor
		if err := rows.Scan(&lm.MonitorId, &lm.Name, &lm.Description, &lm.Subprobe); err != nil {
			return nil, err
		}
		labelMonitors = append(labelMonitors, &lm)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return labelMonitors, nil
}

// TODO(psingh): Perhaps we want to package these up into one function
// and call with something that implements a "Model" interface which has a tableName()
func isExistingLabel(db *sql.DB, id LabelID) (exists bool) {
	if id == 0 {
		return false
	}

	err := db.QueryRow("SELECT EXISTS (SELECT * FROM labels WHERE labelid = ?)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return
}

func (l *Label) Validate(db *sql.DB) (errs []string) {
	if l.Name == "" {
		errs = append(errs, "Name is required")
	}

	for _, lt := range l.Triggers {
		errs = append(errs, lt.Trigger.Validate()...)
	}

	for _, lm := range l.Monitors {
		errs = append(errs, lm.Validate(db)...)
	}
	return
}

func (l *Label) Save(db *sql.DB) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return
	}
	// TODO(psingh): Package up into a helper function for use with everything in the revere package
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if l.LabelId == 0 {
		l.LabelId, err = l.create(tx)
	} else {
		err = l.update(tx)
	}
	if err != nil {
		return
	}

	for _, lt := range l.Triggers {
		err = lt.save(tx, l.LabelId)
		if err != nil {
			return
		}
	}

	for _, lm := range l.Monitors {
		err = lm.save(tx, l.LabelId)
		if err != nil {
			return
		}
	}
	return
}

func (l *Label) create(tx *sql.Tx) (LabelID, error) {
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO labels(%s) VALUES (?, ?, ?)", allLabelFields))
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(nil, l.Name, l.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return LabelID(id), stmt.Close()
}

func (l *Label) update(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		UPDATE labels
		SET name=?, description=?
		WHERE labelid=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(l.Name, l.Description, l.LabelId)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lm *LabelMonitor) Validate(db *sql.DB) (errs []string) {
	if err := validateSubprobe(lm.Subprobe); err != nil {
		errs = append(errs, err.Error())
	}

	if !isExistingMonitor(db, lm.MonitorId) {
		errs = append(errs, fmt.Sprintf("Invalid monitor: %d", lm.MonitorId))
	}
	return
}

func (lm *LabelMonitor) save(tx *sql.Tx, labelId LabelID) (err error) {
	if lm.Create {
		return lm.create(tx, labelId)
	}
	if lm.Delete {
		return lm.delete(tx, labelId)
	}
	return lm.update(tx, labelId)
}

func (lm *LabelMonitor) create(tx *sql.Tx, labelId LabelID) error {
	stmt, err := tx.Prepare(
		`INSERT INTO labels_monitors(labelid, monitorid, subprobe) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(labelId, lm.MonitorId, lm.Subprobe)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lm *LabelMonitor) update(tx *sql.Tx, labelId LabelID) error {
	stmt, err := tx.Prepare(`
		UPDATE labels_monitors
		SET subprobe=?
		WHERE labelid=? AND monitorid=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(lm.Subprobe, labelId, lm.MonitorId)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lm *LabelMonitor) delete(tx *sql.Tx, labelId LabelID) error {
	stmt, err := tx.Prepare(`
		DELETE FROM labels_monitors
		WHERE labelid=? AND monitorid=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(labelId, lm.MonitorId)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lt *LabelTrigger) save(tx *sql.Tx, labelId LabelID) (err error) {
	if lt.Delete {
		return lt.delete(tx)
	}

	newTriggerId, err := lt.Trigger.save(tx)
	if err != nil {
		return err
	}
	if lt.TriggerId == 0 {
		lt.TriggerId = newTriggerId
		err = lt.create(tx, labelId)
	}
	return
}

func (lt *LabelTrigger) create(tx *sql.Tx, labelId LabelID) error {
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO label_triggers(%s) VALUES (?, ?)", allLabelTriggerFields))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(labelId, lt.TriggerId)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lt *LabelTrigger) delete(tx *sql.Tx) error {
	// Trigger delete will cascade to label triggers
	return lt.Trigger.delete(tx)
}
