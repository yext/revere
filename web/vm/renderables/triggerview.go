package renderables

import (
	"github.com/yext/revere/web/vm"
)

type TriggerView struct {
	viewmodel *vm.Trigger
	subs      []Renderable
}

func NewTriggerView(t *vm.Trigger) *TriggerView {
	tv := TriggerView{}
	tv.viewmodel = t
	tv.subs = []Renderable{
		NewTargetView(t.Target),
	}
	return &tv
}

func (tv *TriggerView) name() string {
	return "Trigger"
}

func (tv *TriggerView) template() string {
	return "partials/trigger-view.html"
}

func (tv *TriggerView) data() interface{} {
	return tv.viewmodel
}

func (tv *TriggerView) scripts() []string {
	return nil
}

func (tv *TriggerView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tv *TriggerView) subRenderables() []Renderable {
	return tv.subs
}

func (tv *TriggerView) renderPropagate() (*renderResult, error) {
	return renderPropagate(tv)
}

func (tv *TriggerView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
