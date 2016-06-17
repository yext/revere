package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelTriggerView struct {
	viewmodel *vm.Trigger
	subs      []Renderable
}

func NewLabelTriggerView(t *vm.LabelTrigger) *LabelTriggerView {
	tv := LabelTriggerView{}
	tv.viewmodel = t.Trigger
	tv.subs = []Renderable{
		NewTargetView(t.Trigger.Target),
	}
	return &tv
}

func (tv *LabelTriggerView) name() string {
	return "LabelTrigger"
}

func (tv *LabelTriggerView) template() string {
	return "partials/labeltrigger-view.html"
}

func (tv *LabelTriggerView) data() interface{} {
	return tv.viewmodel
}

func (tv *LabelTriggerView) scripts() []string {
	return nil
}

func (tv *LabelTriggerView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tv *LabelTriggerView) subRenderables() []Renderable {
	return tv.subs
}

func (tv *LabelTriggerView) renderPropagate() (*renderResult, error) {
	return renderPropagate(tv)
}

func (tv *LabelTriggerView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
