package renderables

import (
	"github.com/yext/revere/web/vm"
)

type TriggerEdit struct {
	viewmodel *vm.Trigger
	subs      []Renderable
}

func NewTriggerEdit(t *vm.Trigger) *TriggerEdit {
	te := TriggerEdit{}
	te.viewmodel = t
	te.subs = []Renderable{
		NewTargetEdit(t.Target),
	}
	return &te
}

func (te *TriggerEdit) name() string {
	return "Trigger"
}

func (te *TriggerEdit) template() string {
	return "partials/trigger-edit.html"
}

func (te *TriggerEdit) data() interface{} {
	return te.viewmodel
}

func (te *TriggerEdit) scripts() []string {
	return nil
}

func (te *TriggerEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (te *TriggerEdit) subRenderables() []Renderable {
	return te.subs
}

func (te *TriggerEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(te)
}

func (te *TriggerEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
