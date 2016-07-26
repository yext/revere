package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SilenceView struct {
	silence    *vm.Silence
	subs       []Renderable
	saveStatus string
}

func NewSilenceView(s *vm.Silence, saveStatus []byte) *SilenceView {
	sv := SilenceView{}
	sv.silence = s
	sv.saveStatus = string(saveStatus)

	return &sv
}

func (sv *SilenceView) name() string {
	return "Silence"
}

func (sv *SilenceView) template() string {
	return "silences-view.html"
}

func (sv *SilenceView) data() interface{} {
	return map[string]interface{}{
		"Silence":    sv.silence,
		"SaveStatus": sv.saveStatus,
	}
}

func (sv *SilenceView) scripts() []string {
	return nil
}

func (sv *SilenceView) breadcrumbs() []vm.Breadcrumb {
	return vm.SilencesViewBcs(sv.silence.Id(), sv.silence.MonitorName)
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
