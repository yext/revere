package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SettingsIndex struct {
	settings []*vm.Setting
	subs     []Renderable
}

func NewSettingsIndex(ss []*vm.Setting) *SettingsIndex {
	si := new(SettingsIndex)
	si.settings = ss
	si.subs = make([]Renderable, len(ss))
	for i, s := range ss {
		si.subs[i] = NewSettingEdit(s)
	}
	return si
}

func (si *SettingsIndex) name() string {
	return "SettingsIndex"
}

func (si *SettingsIndex) template() string {
	return "settings-index.html"
}

func (si *SettingsIndex) data() interface{} {
	return si.settings
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
