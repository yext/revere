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

func (s *Setting) Id() int64 {
	return int64(s.Setting.SettingId)
}

func BlankSetting(id settings.SettingTypeId) (*Setting, error) {
	var err error
	viewmodel := new(Setting)
	viewmodel.Setting = new(revere.Setting)
	viewmodel.SettingType, err = settings.SettingTypeById(id)
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
	allTypes := settings.AllSettingTypes()
	var viewmodels []*Setting
	for _, t := range allTypes {
		viewmodelsOfType, err := allSettingsOfType(t, db)
		if err != nil {
			return nil, err
		}

		if len(viewmodelsOfType) == 0 {
			blank, err := BlankSetting(t.Id())
			if err != nil {
				return nil, err
			}
			viewmodels = append(viewmodels, blank)
		}
		viewmodels = append(viewmodels, viewmodelsOfType...)
	}

	return viewmodels, nil
}

func allSettingsOfType(t settings.SettingType, db *sql.DB) ([]*Setting, error) {
	settings, err := revere.LoadSettingsOfType(db, t.Id())
	viewmodels := make([]*Setting, len(settings))
	for i, setting := range settings {
		viewmodels[i], err = newSettingFromModel(setting)
		if err != nil {
			return nil, err
		}
	}
	return viewmodels, nil
}
