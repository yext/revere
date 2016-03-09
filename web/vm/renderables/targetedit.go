package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/web/vm"
)

type TargetEdit struct {
	viewmodel *vm.Target
	subs      []Renderable
}

func NewTargetEdit(t *vm.Target) *TargetEdit {
	te := TargetEdit{}
	te.viewmodel = t
	te.subs = []Renderable{}
	return &te
}

func (te *TargetEdit) name() string {
	return "Target"
}

func (te *TargetEdit) template() string {
	tmpl, ok := te.viewmodel.TargetType().Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for target type %s", te.viewmodel.Target.TargetType().Name()))
	}

	return path.Join(vm.TargetsDir, tmpl)
}

func (te *TargetEdit) data() interface{} {
	return te.viewmodel.Target
}

func (te *TargetEdit) scripts() []string {
	scripts := te.viewmodel.TargetType().Scripts()["edit"]

	return vm.AppendDir(vm.TargetsDir, scripts)
}

func (te *TargetEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (te *TargetEdit) subRenderables() []Renderable {
	return te.subs
}

func (te *TargetEdit) renderPropogate() (*renderResult, error) {
	return renderPropogateImmediate(te)
}

func (te *TargetEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
