package renderables

import (
	"github.com/yext/revere/probe"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

type ProbePreview struct {
	probe probe.VM
}

func NewProbePreview(probe probe.VM) *ProbePreview {
	pp := ProbePreview{}
	pp.probe = probe
	return &pp
}

func (pp *ProbePreview) name() string {
	return "Probe Preview"
}

func (pp *ProbePreview) template() string {
	return "preview.html"
}

func (pp *ProbePreview) data() interface{} {
	return nil
}

func (pp *ProbePreview) scripts() []string {
	return tmpl.AppendDir(probe.ProbesDir, pp.probe.Scripts()["preview"])
}

func (pp *ProbePreview) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pp *ProbePreview) subRenderables() []Renderable {
	return nil
}

func (pp *ProbePreview) renderPropagate() (*renderResult, error) {
	return renderPropagate(pp)
}

func (pp *ProbePreview) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
