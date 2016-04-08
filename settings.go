package revere

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere/settings"
)

type Setting struct {
	Id          uint                   `json:"id"`
	SettingType settings.SettingTypeId `json:"settingType"`
	Setting     string                 `json:"setting"`
	Delete      bool                   `json:"delete,omitempty"`
}

const allSettingFields = `id, settingType, setting`

func (s *Setting) Validate() (errs []string) {
	settingType, err := settings.SettingTypeById(s.SettingType)
	if err != nil {
		errs = append(errs, err.Error())
	}
	setting, err := settingType.Load(s.Setting)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Invalid setting: %s", s.Setting))
	}
	errs = append(errs, setting.Validate()...)
	return
}

func LoadSettingsOfType(db *sql.DB, settingType settings.SettingTypeId) ([]*Setting, error) {
	return LoadSettings(db, fmt.Sprintf("WHERE settingType = %d", settingType))
}

func LoadSettingById(db *sql.DB, id uint) (*Setting, error) {
	results, err := LoadSettings(db, fmt.Sprintf("WHERE id = %d", id))
	if len(results) == 0 {
		return nil, fmt.Errorf("Setting source not found: %d", id)
	}
	return results[0], err
}

func LoadAllSettings(db *sql.DB) ([]*Setting, error) {
	return LoadSettings(db, "")
}

func LoadSettings(db *sql.DB, condition string) (settings []*Setting, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM settings %s", allSettingFields, condition))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ss, err := loadSettingFromRow(rows)
		if err != nil {
			return nil, err
		}
		settings = append(settings, ss)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *Setting) Save(db *sql.DB) (err error) {
	return Transact(db, func(tx *sql.Tx) error {
		if s.isCreate() {
			_, err = s.create(tx)
		} else if s.Delete {
			err = s.delete(tx)
		} else {
			err = s.update(tx)
		}
		return err
	})
}

func (s *Setting) isCreate() bool {
	return s.Id == 0
}

func (s *Setting) create(tx *sql.Tx) (uint, error) {
	stmt, err := tx.Prepare(fmt.Sprintf(`INSERT INTO settings (%s) VALUES (?, ?, ?)`, allSettingFields))
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(s.SettingType, s.Setting)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint(id), err
}

func (s *Setting) update(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE settings SET settingType = ?, setting = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.SettingType, s.Setting, s.Id)
	return err
}

func (s *Setting) delete(tx *sql.Tx) error {
	var stmt *sql.Stmt
	stmt, err := tx.Prepare(`
		DELETE FROM settings
		WHERE id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func loadSettingFromRow(rows *sql.Rows) (*Setting, error) {
	var s Setting
	if err := rows.Scan(&s.Id, &s.SettingType, &s.Setting); err != nil {
		return nil, err
	}
	return &s, nil
}
