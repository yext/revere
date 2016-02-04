package vm

type MonitorEdit struct {
	viewmodel *Monitor
	subs      map[string]Renderable
}

func NewMonitorEdit(m *Monitor) *MonitorEdit {
	mv := MonitorEdit{}
	mv.viewmodel = m
	mv.subs = map[string]Renderable{
	//"Probe":    m.NewProbeEdit(m.Probe),
	//"Triggers": m.TriggersEdit(m.Triggers),
	}
	return &mv
}

func (mv *MonitorEdit) Template() string {
	return "monitors-edit.html"
}

func (mv *MonitorEdit) Data() interface{} {
	return mv.viewmodel
}

func (mv *MonitorEdit) Scripts() []string {
	return []string{
		"revere.js",
		"monitors-edit.js",
		"probes/graphite-preview.js",
		"targets/email.js",
	}
}

func (mv *MonitorEdit) Breadcrumbs() []Breadcrumb {
	return []Breadcrumb{}
}

func (mv *MonitorEdit) SubRenderables() map[string]Renderable {
	return mv.subs
}
