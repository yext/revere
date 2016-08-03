package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/target"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

type TargetView struct {
	viewmodel target.TargetVM
	subs      []Renderable
}

func NewTargetView(t target.TargetVM) *TargetView {
	tv := TargetView{}
	tv.viewmodel = t
	tv.subs = []Renderable{}
	return &tv
}

func (tv *TargetView) name() string {
	return "Target"
}

func (tv *TargetView) template() string {
	tmpl, ok := tv.viewmodel.Type().Templates()["view"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for target type %s", tv.viewmodel.Type().Name()))
	}

	return path.Join(target.TargetsDir, tmpl)
}

func (tv *TargetView) data() interface{} {
	return tv.viewmodel
}

func (tv *TargetView) scripts() []string {
	scripts := tv.viewmodel.Type().Scripts()["view"]

	return tmpl.AppendDir(target.TargetsDir, scripts)
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
