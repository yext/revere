package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/web/vm"
)

type TargetEdit struct {
	viewmodel *vm.Target
	subs      map[string]Renderable
}

func NewTargetEdit(t *vm.Target) *TargetEdit {
	te := TargetEdit{}
	te.viewmodel = t
	te.subs = map[string]Renderable{}
	return &te
}

func (te *TargetEdit) Template() string {
	tmpl, ok := te.viewmodel.TargetType().Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for target type %s", te.viewmodel.Target.TargetType().Name()))
	}

	return path.Join(vm.TargetsDir, tmpl)
}

func (te *TargetEdit) Data() interface{} {
	return te.viewmodel.Target
}

func (te *TargetEdit) Scripts() []string {
	scripts := te.viewmodel.TargetType().Scripts()["edit"]

	return vm.AppendDir(vm.TargetsDir, scripts)
}

func (te *TargetEdit) Breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (te *TargetEdit) SubRenderables() map[string]Renderable {
	return te.subs
}

func (te *TargetEdit) RenderNow() bool {
	return true
}
