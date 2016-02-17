package vm

type MonitorView struct {
	viewmodel *Monitor
	subs      map[string]Renderable
}

func NewMonitorView(m *Monitor) *MonitorView {
	mv := MonitorView{}
	mv.viewmodel = m
	mv.subs = map[string]Renderable{
		"Probe": NewProbeView(m.Probe),
		//"Triggers": NewTriggersView(m.Triggers),
	}
	return &mv
}

func (mv *MonitorView) Template() string {
	return "monitors-view.html"
}

func (mv *MonitorView) Data() interface{} {
	return mv.viewmodel
}

func (mv *MonitorView) Scripts() []string {
	return []string{}
}

func (mv *MonitorView) Breadcrumbs() []Breadcrumb {
	return MonitorViewBcs(mv.viewmodel.Name, mv.viewmodel.Id)
}

func (mv *MonitorView) SubRenderables() map[string]Renderable {
	return mv.subs
}

func (mv *MonitorView) RenderNow() bool {
	return false
}
