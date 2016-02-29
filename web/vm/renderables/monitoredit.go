package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorEdit struct {
	viewmodel *vm.Monitor
	subs      map[string]Renderable
}

func NewMonitorEdit(m *vm.Monitor) *MonitorEdit {
	me := MonitorEdit{}
	me.viewmodel = m
	me.subs = map[string]Renderable{
		"Probe": NewProbeEdit(m.Probe),
		//"Triggers": m.TriggersEdit(m.Triggers),
	}
	return &me
}

func (me *MonitorEdit) Template() string {
	return "monitors-edit.html"
}

func (me *MonitorEdit) Data() interface{} {
	return me.viewmodel
}

func (me *MonitorEdit) Scripts() []string {
	return []string{
		"revere.js",
		"monitors-edit.js",
		"targets/email.js",
	}
}

func (me *MonitorEdit) Breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (me *MonitorEdit) SubRenderables() map[string]Renderable {
	return me.subs
}

func (me *MonitorEdit) RenderNow() bool {
	return false
}
