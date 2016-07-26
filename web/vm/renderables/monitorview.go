package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorView struct {
	monitor    *vm.Monitor
	subs       []Renderable
	saveStatus string
}

func NewMonitorView(m *vm.Monitor, saveStatus []byte) *MonitorView {
	mv := MonitorView{}
	mv.monitor = m
	mv.subs = []Renderable{
		NewProbeView(m.Probe),
		NewMonitorTriggersView(m.Triggers),
		NewMonitorLabelsView(m.Labels),
	}
	mv.saveStatus = string(saveStatus)
	return &mv
}

func (mv *MonitorView) name() string {
	return "Monitor"
}

func (mv *MonitorView) template() string {
	return "monitors-view.html"
}

func (mv *MonitorView) data() interface{} {
	return map[string]interface{}{
		"Monitor":    mv.monitor,
		"SaveStatus": mv.saveStatus,
	}
}

func (mv *MonitorView) scripts() []string {
	return nil
}

func (mv *MonitorView) breadcrumbs() []vm.Breadcrumb {
	return vm.MonitorViewBcs(mv.monitor.Name, mv.monitor.Id())
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
