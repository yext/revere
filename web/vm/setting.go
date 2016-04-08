package vm

import (
	"database/sql"

	"github.com/yext/revere"
	"github.com/yext/revere/settings"
)

const (
	SettingsDir = "settings"
)

type Setting struct {
	*revere.Setting
	Attributes  settings.Setting
	SettingType settings.SettingType
}

func BlankSetting(id int) (*Setting, error) {
	var err error
	viewmodel := new(Setting)
	viewmodel.Setting = new(revere.Setting)
	viewmodel.SettingType, err = settings.SettingTypeById(settings.SettingTypeId(id))
	if err != nil {
		return nil, err
	}
	viewmodel.Attributes, err = viewmodel.SettingType.Load(`{}`)
	if err != nil {
		return nil, err
	}

	return viewmodel, nil
}

func newSettingFromModel(s *revere.Setting) (*Setting, error) {
	viewmodel := new(Setting)
	viewmodel.Setting = s

	settingType, err := settings.SettingTypeById(s.SettingType)
	if err != nil {
		return nil, err
	}

	viewmodel.Attributes, err = settingType.Load(s.Setting)
	if err != nil {
		return nil, err
	}

	viewmodel.SettingType = settingType
	return viewmodel, nil
}

func AllSettings(db *sql.DB) ([]*Setting, error) {
	settings, err := revere.LoadAllSettings(db)
	viewmodels := make([]*Setting, len(settings))
	for i, setting := range settings {
		viewmodels[i], err = newSettingFromModel(setting)
		if err != nil {
			return nil, err
		}
	}
	return viewmodels, nil
}
