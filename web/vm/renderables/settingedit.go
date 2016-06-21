package renderables

import (
	"path"

	"github.com/yext/revere/settings"
	"github.com/yext/revere/web/vm"
)

type SettingEdit struct {
	*settings.VM
}

func NewSettingEdit(s *settings.VM) *SettingEdit {
	se := new(SettingEdit)
	se.VM = s
	return se
}

func (se *SettingEdit) name() string {
	return se.Setting.Type().Name()
}

func (se *SettingEdit) template() string {
	return path.Join(settings.SettingDir, se.Setting.Type().Template())
}

func (se *SettingEdit) data() interface{} {
	return se
}

func (se *SettingEdit) scripts() []string {
	return vm.AppendDir(settings.SettingDir, se.Setting.Type().Scripts())
}

func (se *SettingEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (se *SettingEdit) subRenderables() []Renderable {
	return nil
}

func (se *SettingEdit) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(se)
}

func (se *SettingEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
