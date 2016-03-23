package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SilenceView struct {
	viewmodel *vm.Silence
	subs      []Renderable
}

func NewSilenceView(s *vm.Silence) *SilenceView {
	sv := SilenceView{}
	sv.viewmodel = s

	return &sv
}

func (sv *SilenceView) name() string {
	return "Silence"
}

func (sv *SilenceView) template() string {
	return "silences-view.html"
}

func (sv *SilenceView) data() interface{} {
	return sv.viewmodel
}

func (sv *SilenceView) scripts() []string {
	return nil
}

func (sv *SilenceView) breadcrumbs() []vm.Breadcrumb {
	return vm.SilencesViewBcs(sv.viewmodel.Id, sv.viewmodel.MonitorName)
}

func (sv *SilenceView) subRenderables() []Renderable {
	return nil
}

func (sv *SilenceView) renderPropagate() (*renderResult, error) {
	return renderPropagate(sv)
}

func (sv *SilenceView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
