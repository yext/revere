package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SubprobesIndex struct {
	subprobes []*vm.Subprobe
	monitor   *vm.Monitor
	subs      []Renderable
}

func NewSubprobesIndex(ss []*vm.Subprobe, m *vm.Monitor) *SubprobesIndex {
	ssi := new(SubprobesIndex)
	ssi.subprobes = ss
	ssi.monitor = m

	return ssi
}

func (ssi *SubprobesIndex) name() string {
	return "SubprobesIndex"
}

func (ssi *SubprobesIndex) template() string {
	return "subprobes-index.html"
}

func (ssi *SubprobesIndex) data() interface{} {
	return map[string]interface{}{
		"Subprobes": ssi.subprobes,
		"Monitor":   ssi.monitor,
	}
}

func (ssi *SubprobesIndex) scripts() []string {
	return []string{
		"entered-states.js",
		"subprobes-index.js",
	}
}

func (ssi *SubprobesIndex) breadcrumbs() []vm.Breadcrumb {
	return vm.SubprobeIndexBcs(ssi.monitor.Name, ssi.monitor.Id())
}

func (ssi *SubprobesIndex) subRenderables() []Renderable {
	return nil
}

func (ssi *SubprobesIndex) renderPropagate() (*renderResult, error) {
	return renderPropagate(ssi)
}

func (ssi *SubprobesIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
