package vm

type MonitorEdit struct {
	viewmodel *Monitor
	subs      map[string]Renderable
}

func NewMonitorEdit(m *Monitor) *MonitorEdit {
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
		"probes/graphite-preview.js",
		"targets/email.js",
	}
}

func (me *MonitorEdit) Breadcrumbs() []Breadcrumb {
	return []Breadcrumb{}
}

func (me *MonitorEdit) SubRenderables() map[string]Renderable {
	return me.subs
}

func (me *MonitorEdit) RenderNow() bool {
	return false
}
