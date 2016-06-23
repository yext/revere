package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelTriggerEdit struct {
	labelTrigger *vm.LabelTrigger
	subs         []Renderable
}

func NewLabelTriggerEdit(t *vm.LabelTrigger) *LabelTriggerEdit {
	te := LabelTriggerEdit{}
	te.labelTrigger = t
	te.subs = []Renderable{
		NewTargetEdit(t.Trigger.Target),
	}
	return &te
}

func (te *LabelTriggerEdit) name() string {
	return "Trigger"
}

func (te *LabelTriggerEdit) template() string {
	return "partials/label-trigger-edit.html"
}

func (te *LabelTriggerEdit) data() interface{} {
	return te.labelTrigger
}

func (te *LabelTriggerEdit) scripts() []string {
	return []string{
		"trigger-edit.js",
	}
}

func (te *LabelTriggerEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (te *LabelTriggerEdit) subRenderables() []Renderable {
	return te.subs
}

func (te *LabelTriggerEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(te)
}

func (te *LabelTriggerEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
