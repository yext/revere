package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelTriggerEdit struct {
	viewmodel *vm.Trigger
	subs      []Renderable
}

func NewLabelTriggerEdit(t *vm.LabelTrigger) *LabelTriggerEdit {
	te := LabelTriggerEdit{}
	te.viewmodel = t.Trigger
	te.subs = []Renderable{
		NewTargetEdit(t.Trigger.Target),
	}
	return &te
}

func (te *LabelTriggerEdit) name() string {
	return "LabelTrigger"
}

func (te *LabelTriggerEdit) template() string {
	return "partials/labeltrigger-edit.html"
}

func (te *LabelTriggerEdit) data() interface{} {
	return te.viewmodel
}

func (te *LabelTriggerEdit) scripts() []string {
	return nil
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
