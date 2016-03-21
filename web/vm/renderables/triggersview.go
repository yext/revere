package renderables

import (
	"github.com/yext/revere/web/vm"
)

type TriggersView struct {
	viewmodel []*vm.Trigger
	subs      []Renderable
}

func NewTriggersView(ts []*vm.Trigger) *TriggersView {
	tsv := TriggersView{}
	for _, trigger := range ts {
		tsv.subs = append(tsv.subs, NewTriggerView(trigger))
	}

	return &tsv
}

func (tsv *TriggersView) name() string {
	return "Triggers"
}

func (tsv *TriggersView) template() string {
	return "partials/triggers-view.html"
}

func (tsv *TriggersView) data() interface{} {
	return tsv.viewmodel
}

func (tsv *TriggersView) scripts() []string {
	return nil
}

func (tsv *TriggersView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tsv *TriggersView) subRenderables() []Renderable {
	return tsv.subs
}

func (tsv *TriggersView) renderPropagate() (*renderResult, error) {
	return renderPropagate(tsv)
}

func (tsv *TriggersView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
