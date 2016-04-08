package renderables

import (
	"path"

	"github.com/yext/revere/web/vm"
)

type SettingEdit struct {
	*vm.Setting
}

func NewSettingEdit(s *vm.Setting) *SettingEdit {
	se := new(SettingEdit)
	se.Setting = s
	return se
}

func (se *SettingEdit) name() string {
	return se.SettingType.Name()
}

func (se *SettingEdit) template() string {
	return path.Join(vm.SettingsDir, se.SettingType.Template())
}

func (se *SettingEdit) data() interface{} {
	return se
}

func (se *SettingEdit) scripts() []string {
	return vm.AppendDir(vm.SettingsDir, se.SettingType.Scripts())
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
