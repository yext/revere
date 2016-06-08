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
)

type SettingTypeId int16

type SettingType interface {
	Id() SettingTypeId
	Name() string
	loadFromParams(ds string) (Setting, error)
	loadFromDb(ds string) (Setting, error)
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
	types       = make(map[SettingTypeId]SettingType)
)

func Default() (Setting, error) {
	s, err := defaultType.blank()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return s, nil
}

func LoadFromParams(id SettingTypeId, sParams string) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return sType.loadFromParams(sParams)
}

func LoadFromDb(id SettingTypeId, sJson string) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return sType.loadFromDb(sJson)
}

func Blank(id SettingTypeId) (Setting, error) {
	sType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return sType.blank()
}

func getType(id SettingTypeId) (SettingType, error) {
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
