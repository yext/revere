package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

type ProbeEdit struct {
	probe probes.Probe
	subs  []Renderable
}

func NewProbeEdit(p probes.Probe) *ProbeEdit {
	pe := ProbeEdit{}
	pe.probe = p
	pe.subs = []Renderable{}
	return &pe
}

func (pe *ProbeEdit) name() string {
	return "Probe"
}

func (pe *ProbeEdit) template() string {
	tmpl, ok := pe.probe.Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pe.probe.Name()))
	}

	return path.Join(probes.ProbesDir, tmpl)
}

func (pe *ProbeEdit) data() interface{} {
	return map[string]interface{}{
		"Probe": pe.probe,
	}
}

func (pe *ProbeEdit) scripts() []string {
	scripts := pe.probe.Scripts()["edit"]

	return tmpl.AppendDir(probes.ProbesDir, scripts)
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
