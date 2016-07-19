package renderables

import (
	"fmt"
	"regexp"

	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/vm"
)

type SubprobeView struct {
	subprobe *vm.Subprobe
	readings []*vm.Reading
	probe    probes.Probe
	subs     []Renderable
}

func NewSubprobeView(p probes.Probe, s *vm.Subprobe, rs []*vm.Reading) *SubprobeView {
	sv := SubprobeView{}
	sv.subprobe = s
	sv.probe = p
	sv.readings = rs
	pp := NewProbePreview(p)
	sv.subs = []Renderable{pp}
	return &sv
}

func (sv *SubprobeView) name() string {
	return "Subprobe"
}

func (sv *SubprobeView) template() string {
	return "subprobes-view.html"
}

func (sv *SubprobeView) data() interface{} {
	return map[string]interface{}{
		"EscapedSubprobeName": fmt.Sprintf("^%s$", regexp.QuoteMeta(sv.subprobe.Name)),
		"Subprobe":            sv.subprobe,
		"Readings":            sv.readings,
		"PreviewParams":       sv.probe.SerializeForFrontend(),
	}
}

func (sv *SubprobeView) scripts() []string {
	return []string{
		"subprobes-view.js",
	}
}

func (sv *SubprobeView) breadcrumbs() []vm.Breadcrumb {
	return vm.SubprobeViewBcs(sv.subprobe)
}

func (sv *SubprobeView) subRenderables() []Renderable {
	return sv.subs
}

func (sv *SubprobeView) renderPropagate() (*renderResult, error) {
	return renderPropagate(sv)
}

func (sv *SubprobeView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
