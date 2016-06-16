package renderables

import (
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
		"Subprobe": sv.subprobe,
		"Readings": sv.readings,
	}
}

func (sv *SubprobeView) scripts() []string {
	return nil
}

func (sv *SubprobeView) breadcrumbs() []vm.Breadcrumb {
	return vm.SubprobeViewBcs(sv.subprobe.Subprobe)
}

func (sv *SubprobeView) subRenderables() []Renderable {
	return nil
}

func (sv *SubprobeView) renderPropagate() (*renderResult, error) {
	return renderPropagate(sv)
}

func (sv *SubprobeView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
