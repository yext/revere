/*
   Package settings provides a framework for adding custom settings to your-revere.com/settings/ page.
   Add your setting to this package and have it implement the Setting interface.
   settingTemplates contains the templates in settingTemplateDir, so they are available for use in your setting.
   Don't forget to register your setting with registerSetting(Setting)
*/
package settings

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type VM struct {
	Setting
	SettingParams string
	SettingType   db.SettingType
	SettingID     db.SettingID
}

type SettingType interface {
	Id() db.SettingType
	Name() string
	loadFromParams(ds string) (Setting, error)
	loadFromDB(ds string) (Setting, error)
	blank() (Setting, error)
	Template() string
	Scripts() []string
}

type Setting interface {
	SettingType
	Serialize() (string, error)
	Type() SettingType
	Validate() []string
}

const (
	SettingDir = "settings"
)

var (
	defaultType = OutgoingEmail{}
	types       = make(map[db.SettingType]SettingType)
)

func Default() (Setting, error) {
	s, err := defaultType.blank()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return s, nil
}

func LoadFromParams(id db.SettingType, sParams string) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return sType.loadFromParams(sParams)
}

func LoadFromDB(id db.SettingType, sJson string) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return sType.loadFromDB(sJson)
}

func Blank(id db.SettingType) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return sType.blank()
}

func getType(id db.SettingType) (SettingType, error) {
	if s, ok := types[id]; !ok {
		return nil, errors.Errorf("Invalid setting type with id: %d", id)
	} else {
		return s, nil
	}
}

func addSettingType(s SettingType) {
	if _, exists := types[s.Id()]; !exists {
		types[s.Id()] = s
	} else {
		panic(fmt.Sprintf("A setting type with id %d already exists", s.Id()))
	}
}

func AllTypes() (sts []SettingType) {
	for _, t := range types {
		sts = append(sts, t)
	}
	return sts
}

func All(DB *db.DB) ([]*VM, error) {
	settings, err := DB.LoadSettings()
	if err != nil {
		return nil, err
	}

	ss := make([]*VM, len(settings))
	for i, setting := range settings {
		ss[i], err = newVM(setting)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return ss, nil
}

func newVM(s *db.Setting) (*VM, error) {
	setting, err := LoadFromDB(s.SettingType, s.Setting)
	if err != nil {
		return &VM{}, errors.Trace(err)
	}

	return &VM{
		Setting:     setting,
		SettingID:   s.SettingID,
		SettingType: s.SettingType,
	}, nil
}

func (vm *VM) Save(tx *db.Tx) error {
	var err error
	vm.Setting, err = LoadFromParams(vm.SettingType, vm.SettingParams)
	if err != nil {
		return errors.Trace(err)
	}

	settingJSON, err := vm.Setting.Serialize()
	if err != nil {
		return errors.Trace(err)
	}

	setting := &db.Setting{
		SettingID:   vm.SettingID,
		SettingType: vm.SettingType,
		Setting:     settingJSON,
	}

	err = tx.UpdateSetting(setting)

	return errors.Trace(err)
}

func (vm *VM) Validate() (errs []string) {
	var err error
	vm.Setting, err = LoadFromParams(vm.SettingType, vm.SettingParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load setting %s", vm.SettingParams))
		return errs
	}

	return vm.Setting.Validate()
}
