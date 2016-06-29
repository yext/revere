package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

type ProbeView struct {
	viewmodel probes.Probe
	subs      []Renderable
}

func NewProbeView(p probes.Probe) *ProbeView {
	pv := ProbeView{}
	pv.viewmodel = p
	pv.subs = []Renderable{}
	return &pv
}

func (pv *ProbeView) name() string {
	return "Probe"
}

func (pv *ProbeView) template() string {
	tmpl, ok := pv.viewmodel.Type().Templates()["view"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pv.viewmodel.Type().Name()))
	}

	return path.Join(probes.ProbesDir, tmpl)
}

func (pv *ProbeView) data() interface{} {
	return pv.viewmodel
}

func (pv *ProbeView) scripts() []string {
	scripts := pv.viewmodel.Type().Scripts()["view"]

	return tmpl.AppendDir(probes.ProbesDir, scripts)
}

func (pv *ProbeView) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pv *ProbeView) subRenderables() []Renderable {
	return pv.subs
}

func (pv *ProbeView) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(pv)
}

func (pv *ProbeView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
