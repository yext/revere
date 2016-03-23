package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorEdit struct {
	viewmodel *vm.Monitor
	subs      []Renderable
}

func NewMonitorEdit(m *vm.Monitor) *MonitorEdit {
	me := MonitorEdit{}
	me.viewmodel = m
	me.subs = []Renderable{
		NewProbeEdit(m.Probe),
		NewTriggersEdit(m.Triggers),
	}
	return &me
}

func (me *MonitorEdit) name() string {
	return "Monitor"
}

func (me *MonitorEdit) template() string {
	return "monitors-edit.html"
}

func (me *MonitorEdit) data() interface{} {
	return me.viewmodel
}

func (me *MonitorEdit) scripts() []string {
	return []string{
		"monitors-edit.js",
	}
}

func (me *MonitorEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (me *MonitorEdit) subRenderables() []Renderable {
	return me.subs
}

func (me *MonitorEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(me)
}

func (me *MonitorEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
