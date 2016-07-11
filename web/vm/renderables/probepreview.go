package renderables

import (
	"github.com/yext/revere/web/vm"
)

type ProbePreview struct {
	script string
	subs   []Renderable
}

func NewProbePreview(script string) *ProbePreview {
	pp := ProbePreview{}
	pp.script = script
	pp.subs = []Renderable{}
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
	return []string{pp.script}
}

func (pp *ProbePreview) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pp *ProbePreview) subRenderables() []Renderable {
	return pp.subs
}

func (pp *ProbePreview) renderPropagate() (*renderResult, error) {
	return renderPropagate(pp)
}

func (pp *ProbePreview) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
