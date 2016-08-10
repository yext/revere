/*
Package settings provides a framework for adding custom settings to
your-revere.com/settings/ page. Add your setting to this package and have it
implement the Setting interface. Don't forget to register your setting with
addSettingType(SettingType).
*/
package setting

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

// The VM struct is practically identical in purpose to its counterparts in the
// vm package, as it represents the intermediate structure between Revere's DB
// representation of the Setting and its front end representation of the
// Setting.
type VM struct {
	Setting
	SettingParams string
	SettingType   db.SettingType
	SettingID     db.SettingID
}

// SettingType and Setting define a common display abstraction for all
// settings.
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
		return nil, errors.Trace(err)
	}

	return sType.loadFromParams(sParams)
}

func LoadFromDB(id db.SettingType, sJson string) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return sType.loadFromDB(sJson)
}

func Blank(id db.SettingType) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
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

func addType(s SettingType) {
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
		return nil, errors.Trace(err)
	}

	blankTypes := make(map[SettingType]struct{})
	for _, t := range types {
		blankTypes[t] = struct{}{}
	}

	ss := make([]*VM, len(settings))
	for i, setting := range settings {
		ss[i], err = newVM(setting)
		if err != nil {
			return nil, errors.Trace(err)
		}

		delete(blankTypes, ss[i].Setting.Type())
	}

	for bt, _ := range blankTypes {
		blankSetting, err := bt.blank()
		if err != nil {
			return nil, errors.Trace(err)
		}
		ss = append(ss, &VM{
			Setting:     blankSetting,
			SettingType: bt.Id(),
		})
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

func (vm *VM) Id() int64 {
	return int64(vm.SettingID)
}

func (*VM) ComponentName() string {
	return "Setting"
}

func (vm *VM) IsCreate() bool {
	return vm.Id() == 0
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

	if vm.IsCreate() {
		var id db.SettingID
		id, err = tx.CreateSetting(setting)
		setting.SettingID = id
	} else {
		err = tx.UpdateSetting(setting)
	}

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
