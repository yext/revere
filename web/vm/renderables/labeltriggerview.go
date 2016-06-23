package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelTriggerView struct {
	labelTrigger *vm.LabelTrigger
	subs         []Renderable
}

func NewLabelTriggerView(t *vm.LabelTrigger) *LabelTriggerView {
	tv := LabelTriggerView{}
	tv.labelTrigger = t
	tv.subs = []Renderable{
		NewTargetView(t.Trigger.Target),
	}
	return &tv
}

func (tv *LabelTriggerView) name() string {
	return "Trigger"
}

func (tv *LabelTriggerView) template() string {
	return "partials/trigger-view.html"
}

func (tv *LabelTriggerView) data() interface{} {
	return tv.labelTrigger
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
