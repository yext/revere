package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/target"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

type TargetEdit struct {
	viewmodel target.TargetVM
	subs      []Renderable
}

func NewTargetEdit(t target.TargetVM) *TargetEdit {
	te := TargetEdit{}
	te.viewmodel = t
	te.subs = []Renderable{}
	return &te
}

func (te *TargetEdit) name() string {
	return "Target"
}

func (te *TargetEdit) template() string {
	tmpl, ok := te.viewmodel.Type().Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for target type %s", te.viewmodel.Type().Name()))
	}

	return path.Join(target.TargetsDir, tmpl)
}

func (te *TargetEdit) data() interface{} {
	return te.viewmodel
}

func (te *TargetEdit) scripts() []string {
	scripts := te.viewmodel.Type().Scripts()["edit"]

	return tmpl.AppendDir(target.TargetsDir, scripts)
}

func (te *TargetEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (te *TargetEdit) subRenderables() []Renderable {
	return te.subs
}

func (te *TargetEdit) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(te)
}

func (te *TargetEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
