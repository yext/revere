package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/vm"
)

type ProbeEdit struct {
	viewmodel *probes.Probe
	subs      []Renderable
}

func NewProbeEdit(p *probes.Probe) *ProbeEdit {
	pe := ProbeEdit{}
	pe.viewmodel = p
	pe.subs = []Renderable{}
	return &pe
}

func (pe *ProbeEdit) name() string {
	return "Probe"
}

func (pe *ProbeEdit) template() string {
	tmpl, ok := (*pe.viewmodel).Type().Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", (*pe.viewmodel).Type().Name()))
	}

	return path.Join(probes.ProbesDir, tmpl)
}

func (pe *ProbeEdit) data() interface{} {
	return pe.viewmodel
}

func (pe *ProbeEdit) scripts() []string {
	scripts := (*pe.viewmodel).Type().Scripts()["edit"]

	return vm.AppendDir(probes.ProbesDir, scripts)
}

func (pe *ProbeEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pe *ProbeEdit) subRenderables() []Renderable {
	return pe.subs
}

func (pe *ProbeEdit) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(pe)
}

func (pe *ProbeEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
