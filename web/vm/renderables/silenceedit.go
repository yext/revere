package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SilenceEdit struct {
	viewmodel *vm.Silence
	subs      []Renderable
}

func NewSilenceEdit(s *vm.Silence) *SilenceEdit {
	se := SilenceEdit{}
	se.viewmodel = s

	return &se
}

func (se *SilenceEdit) name() string {
	return "Silence"
}

func (se *SilenceEdit) template() string {
	return "silences-edit.html"
}

func (se *SilenceEdit) data() interface{} {
	return se.viewmodel
}

func (se *SilenceEdit) scripts() []string {
	return []string{
		"silences.js",
	}
}

func (se *SilenceEdit) breadcrumbs() []vm.Breadcrumb {
	return vm.SilencesViewBcs(se.viewmodel.Id, se.viewmodel.MonitorName)
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
