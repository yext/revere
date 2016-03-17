package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorView struct {
	viewmodel *vm.Monitor
	subs      []Renderable
}

func NewMonitorView(m *vm.Monitor) *MonitorView {
	mv := MonitorView{}
	mv.viewmodel = m
	mv.subs = []Renderable{
		NewProbeView(m.Probe),
		//NewTriggersView(m.Triggers),
	}
	return &mv
}

func (mv *MonitorView) name() string {
	return "Monitor"
}

func (mv *MonitorView) template() string {
	return "monitors-view.html"
}

func (mv *MonitorView) data() interface{} {
	return mv.viewmodel
}

func (mv *MonitorView) scripts() []string {
	return nil
}

func (mv *MonitorView) breadcrumbs() []vm.Breadcrumb {
	return vm.MonitorViewBcs(mv.viewmodel.Name, mv.viewmodel.Id)
}

func (mv *MonitorView) subRenderables() []Renderable {
	return mv.subs
}

func (mv *MonitorView) renderPropagate() (*renderResult, error) {
	return renderPropagate(mv)
}

func (mv *MonitorView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
