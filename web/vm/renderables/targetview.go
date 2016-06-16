package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/targets"
	"github.com/yext/revere/web/vm"
)

type TargetView struct {
	viewmodel *targets.VM
	subs      []Renderable
}

func NewTargetView(t *vm.Target) *TargetView {
	tv := TargetView{}
	tv.viewmodel = t
	tv.subs = []Renderable{}
	return &tv
}

func (tv *TargetView) name() string {
	return "Target"
}

func (tv *TargetView) template() string {
	tmpl, ok := tv.viewmodel.TargetType().Templates()["view"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for target type %s", tv.viewmodel.Target.TargetType().Name()))
	}

	return path.Join(vm.TargetsDir, tmpl)
}

func (tv *TargetView) data() interface{} {
	return tv.viewmodel.Target
}

func (tv *TargetView) scripts() []string {
	scripts := tv.viewmodel.TargetType().Scripts()["view"]

	return vm.AppendDir(vm.TargetsDir, scripts)
}

func (tv *TargetView) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (tv *TargetView) subRenderables() []Renderable {
	return tv.subs
}

func (tv *TargetView) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(tv)
}

func (tv *TargetView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
