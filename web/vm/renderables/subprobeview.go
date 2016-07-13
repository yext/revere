package renderables

import (
	"fmt"
	"path"
	"regexp"

	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/vm"
)

type SubprobeView struct {
	subprobe *vm.Subprobe
	readings []*vm.Reading
	subs     []Renderable
}

func NewSubprobeView(s *vm.Subprobe, rs []*vm.Reading) *SubprobeView {
	sv := SubprobeView{}
	sv.subprobe = s
	sv.readings = rs
	pp := NewProbePreview(s.Probe.PreviewScript())
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
		"PreviewParams":       sv.subprobe.Probe.PreviewParams(),
	}
}

func (sv *SubprobeView) scripts() []string {
	return []string{
		path.Join(probes.ProbesDir, sv.subprobe.Probe.PreviewScript()),
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
