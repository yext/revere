package renderables

import (
	"github.com/yext/revere/settings"
	"github.com/yext/revere/web/vm"
)

type SettingsIndex struct {
	settings   []*settings.VM
	subs       []Renderable
	saveStatus string
}

func NewSettingsIndex(ss []*settings.VM, saveStatus []byte) *SettingsIndex {
	si := new(SettingsIndex)
	si.settings = ss
	si.subs = make([]Renderable, len(ss))
	for i, s := range ss {
		si.subs[i] = NewSettingEdit(s)
	}
	si.saveStatus = string(saveStatus)
	return si
}

func (si *SettingsIndex) name() string {
	return "SettingsIndex"
}

func (si *SettingsIndex) template() string {
	return "settings-index.html"
}

func (si *SettingsIndex) data() interface{} {
	return map[string]interface{}{
		"Settings":   si.settings,
		"SaveStatus": si.saveStatus,
	}
}

func (si *SettingsIndex) scripts() []string {
	return []string{
		"settings.js",
	}
}

func (si *SettingsIndex) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (si *SettingsIndex) subRenderables() []Renderable {
	return si.subs
}

func (si *SettingsIndex) renderPropagate() (*renderResult, error) {
	return renderPropagate(si)
}

func (si *SettingsIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
