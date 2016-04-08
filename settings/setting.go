/*
   Package settings provides a framework for adding custom settings to your-revere.com/settings/ page.
   Add your setting to this package and have it implement the Setting interface.
   settingTemplates contains the templates in settingTemplateDir, so they are available for use in your setting.
   Don't forget to register your setting with registerSetting(Setting)
*/
package settings

import "fmt"

type SettingTypeId int

type SettingType interface {
	Id() SettingTypeId
	Name() string
	Template() string
	Scripts() []string
	Load(string) (Setting, error)
	LoadDefault() Setting
}

type Setting interface {
	Validate() []string
	SettingType() SettingType
}

var (
	defaultSettingTypeId = OutgoingEmail{}.Id()
	settingTypes         = make(map[SettingTypeId]SettingType)
)

func SettingTypeById(settingType SettingTypeId) (SettingType, error) {
	if s, ok := settingTypes[settingType]; !ok {
		return nil, fmt.Errorf("Invalid setting type with id: %d", settingType)
	} else {
		return s, nil
	}
}

func addSettingType(s SettingType) {
	if _, exists := settingTypes[s.Id()]; !exists {
		settingTypes[s.Id()] = s
	} else {
		panic(fmt.Sprintf("A setting with id %d already exists", s.Id()))
	}
}

func AllSettingTypes() (types []*SettingType) {
	types = make([]*SettingType, len(settingTypes))
	i := 0
	for _, s := range settingTypes {
		types[i] = &s
		i++
	}
	return
}

func DefaultSettingType() SettingType {
	return settingTypes[defaultSettingTypeId]
}
