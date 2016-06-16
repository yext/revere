package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SilenceEdit struct {
	silence  *vm.Silence
	monitors []*vm.Monitor
	subs     []Renderable
}

func NewSilenceEdit(s *vm.Silence, ms []*vm.Monitor) *SilenceEdit {
	se := SilenceEdit{}
	se.silence = s
	se.monitors = ms

	return &se
}

func (se *SilenceEdit) name() string {
	return "Silence"
}

func (se *SilenceEdit) template() string {
	return "silences-edit.html"
}

func (se *SilenceEdit) data() interface{} {
	return map[string]interface{}{
		"Silence":  se.silence,
		"Monitors": se.monitors,
	}
}

func (se *SilenceEdit) scripts() []string {
	return []string{
		"silences.js",
	}
}

func (se *SilenceEdit) breadcrumbs() []vm.Breadcrumb {
	return vm.SilencesViewBcs(se.silence.Id(), se.silence.MonitorName)
}

func (se *SilenceEdit) subRenderables() []Renderable {
	return nil
}

func (se *SilenceEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(se)
}

func (se *SilenceEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
