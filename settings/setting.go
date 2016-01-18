/*
   Package settings provides a framework for adding custom settings to your-revere.com/settings/ page.
   Add your setting to this package and have it implement the Setting interface.
   settingTemplates contains the templates in settingTemplateDir, so they are available for use in your setting.
   Don't forget to register your setting with registerSetting(Setting)
*/
package settings

import (
	"fmt"
	"html/template"

	"github.com/yext/revere/web/tmpl"
)

type Setting interface {
	Id() int
	Load() error
	Save(jsonObject string) error
	Render() (template.HTML, error)
}

const settingTemplateDir = "web/views/settings/"

var (
	settingTemplates map[string]*template.Template
	settings         map[int]Setting = make(map[int]Setting)
)

func init() {
	funcMap := template.FuncMap{}
	// Fetch all setting templates
	settingTemplates = tmpl.InitTemplates(settingTemplateDir, funcMap)
}

func registerSetting(s Setting) {
	_, exists := settings[s.Id()]
	if exists {
		panic(fmt.Sprintf("A setting with id %d already exists", s.Id()))
	}
	settings[s.Id()] = s
}

// GetSetting returns setting with the provided
// id and a boolean for if it's found or not
func GetSetting(id int) (Setting, bool) {
	setting, ok := settings[id]
	return setting, ok
}

// Returns all the settings after calling Load() on each
func GetAllLoadedSettings() (loadedSettings []Setting) {
	for _, v := range settings {
		v.Load()
		loadedSettings = append(loadedSettings, v)
	}
	return
}

// Returns all the settings which may or may not be loaded.
func GetAllUnloadedSettings() (unloadedSettings []Setting) {
	for _, v := range settings {
		unloadedSettings = append(unloadedSettings, v)
	}
	return
}
