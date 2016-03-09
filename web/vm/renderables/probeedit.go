package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/web/vm"
)

type ProbeEdit struct {
	viewmodel *vm.Probe
	subs      []Renderable
}

func NewProbeEdit(p *vm.Probe) *ProbeEdit {
	pe := ProbeEdit{}
	pe.viewmodel = p
	pe.subs = []Renderable{}
	return &pe
}

func (pe *ProbeEdit) name() string {
	return "Probe"
}

func (pe *ProbeEdit) template() string {
	tmpl, ok := pe.viewmodel.ProbeType().Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pe.viewmodel.Probe.ProbeType().Name()))
	}

	return path.Join(vm.ProbesDir, tmpl)
}

func (pe *ProbeEdit) data() interface{} {
	return pe.viewmodel.Probe
}

func (pe *ProbeEdit) scripts() []string {
	scripts := pe.viewmodel.ProbeType().Scripts()["edit"]

	return vm.AppendDir(vm.ProbesDir, scripts)
}

func (pe *ProbeEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pe *ProbeEdit) subRenderables() []Renderable {
	return pe.subs
}

func (pe *ProbeEdit) renderPropogate() (*renderResult, error) {
	return renderPropogateImmediate(pe)
}

func (pe *ProbeEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
