package revere

import (
	"database/sql"
	"fmt"
)

type Label struct {
	Id          uint            `json:"id,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Triggers    []*LabelTrigger `json:"triggers"`
	Monitors    []*LabelMonitor `json:"monitors,omitempty"`
}

type LabelMonitor struct {
	Subprobes string `json:"subprobes"`
	Monitor
	Create bool `json:"create"`
	Delete bool `json:"delete"`
}

type LabelTrigger struct {
	Trigger
	Delete bool `json:"delete"`
}

const (
	allLabelFields        = "id, name, description"
	allLabelMonitorFields = "m.Id, m.Name, m.Description, lm.Subprobes"
	allLabelTriggerFields = "label_id, trigger_id"
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

func LoadLabel(db *sql.DB, id uint) (l *Label, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM labels WHERE id = %d", allLabelFields, id))
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
	if err := rows.Scan(&l.Id, &l.Name, &l.Description); err != nil {
		return nil, err
	}

	return &l, nil
}

func LoadLabelMonitors(db *sql.DB, labelId uint) ([]*LabelMonitor, error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM monitors m
		JOIN labels_monitors lm on m.id=lm.monitor_id
		WHERE lm.label_id = %d
	`, allLabelMonitorFields, labelId))
	if err != nil {
		return nil, err
	}

	labelMonitors := make([]*LabelMonitor, 0)
	for rows.Next() {
		var lm LabelMonitor
		if err := rows.Scan(&lm.Id, &lm.Name, &lm.Description, &lm.Subprobes); err != nil {
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

	if l.Id == 0 {
		l.Id, err = l.create(tx)
	} else {
		err = l.update(tx)
	}
	if err != nil {
		return
	}

	for _, lt := range l.Triggers {
		err = lt.save(tx, l.Id)
		if err != nil {
			return
		}
	}

	for _, lm := range l.Monitors {
		err = lm.save(tx, l.Id)
		if err != nil {
			return
		}
	}
	return
}

func (l *Label) create(tx *sql.Tx) (uint, error) {
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
	return uint(id), stmt.Close()
}

func (l *Label) update(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		UPDATE labels
		SET name=?, description=?
		WHERE id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(l.Name, l.Description, l.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lm *LabelMonitor) Validate(db *sql.DB) (errs []string) {
	if err := validateSubprobes(lm.Subprobes); err != nil {
		errs = append(errs, err.Error())
	}

	if !isExistingMonitor(db, lm.Id) {
		errs = append(errs, fmt.Sprintf("Invalid monitor: %d", lm.Id))
	}
	return
}

func (lm *LabelMonitor) save(tx *sql.Tx, labelId uint) (err error) {
	if lm.Create {
		return lm.create(tx, labelId)
	}
	if lm.Delete {
		return lm.delete(tx, labelId)
	}
	return lm.update(tx, labelId)
}

func (lm *LabelMonitor) create(tx *sql.Tx, labelId uint) error {
	stmt, err := tx.Prepare(
		`INSERT INTO labels_monitors(label_id, monitor_id, subprobes) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(labelId, lm.Id, lm.Subprobes)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lm *LabelMonitor) update(tx *sql.Tx, labelId uint) error {
	stmt, err := tx.Prepare(`
		UPDATE labels_monitors
		SET subprobes=?
		WHERE label_id=? AND monitor_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(lm.Subprobes, labelId, lm.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lm *LabelMonitor) delete(tx *sql.Tx, labelId uint) error {
	stmt, err := tx.Prepare(`
		DELETE FROM labels_monitors
		WHERE label_id=? AND monitor_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(labelId, lm.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lt *LabelTrigger) save(tx *sql.Tx, labelId uint) (err error) {
	newTriggerId, err := lt.Trigger.save(tx)
	if err != nil {
		return err
	}

	if lt.Id == 0 {
		lt.Id = newTriggerId
		err = lt.create(tx, labelId)
	}
	return
}

func (lt *LabelTrigger) create(tx *sql.Tx, labelId uint) error {
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO label_triggers(%s) VALUES (?, ?)", allLabelTriggerFields))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(labelId, lt.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func (lt *LabelTrigger) delete(tx *sql.Tx) error {
	// Trigger delete will cascade to label triggers
	return lt.Trigger.delete(tx)
}
