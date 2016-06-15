package db

import (
	"database/sql"

	"github.com/juju/errors"
)

type SettingID int32
type SettingType int16

type Setting struct {
	SettingID   SettingID
	SettingType SettingType
	Setting     string
}

func (db *DB) LoadSettings() (settings []*Setting, err error) {
	q := "SELECT * FROM pfx_settings ORDER BY settingtype"
	err = db.Select(&settings, cq(db, q))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return settings, nil
}

func (db *DB) LoadSettingsOfType(settingType SettingTypeId) (settings []*Setting, err error) {
	q := "SELECT * FROM pfx_settings WHERE settingtype = ?"
	err = db.Select(&settings, cq(db, q), settingType)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return settings, nil
}

func (db *DB) LoadSettingById(id SettingID) (*Setting, error) {
	var s Setting
	q := "SELECT * FROM pfx_settings WHERE settingid = ?"

	if err := db.Get(&s, cq(db, q), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &s, nil
}

func (tx *Tx) CreateSetting(s *Setting) error {
	q := `INSERT INTO pfx_settings (settingid, settingtype, setting) 
		VALUES (:settingid, :settingtype, :setting)`
	_, err := tx.NamedExec(cq(tx, q), *s)
	return errors.Trace(err)
}

func (tx *Tx) UpdateSetting(s *Setting) error {
	q := `UPDATE pfx_settings SET settingtype=:settingtype, setting=:setting 
		WHERE settingid=:settingid`
	_, err := tx.NamedExec(cq(tx, q), *s)
	return errors.Trace(err)
}

func (tx *Tx) DeleteSetting(s *Setting) error {
	q := `DELETE FROM pfx_settings WHERE settingid=:settingid`
	_, err := tx.NamedExec(cq(tx, q), *s)
	return errors.Trace(err)
}
