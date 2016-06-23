package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelTriggersView struct {
	subs []Renderable
}

func NewLabelTriggersView(ts []*vm.LabelTrigger) *LabelTriggersView {
	tsv := LabelTriggersView{}
	for _, labeltrigger := range ts {
		tsv.subs = append(tsv.subs, NewLabelTriggerView(labeltrigger))
	}

	return &tsv
}

func (tsv *LabelTriggersView) name() string {
	return "Triggers"
}

func (tsv *LabelTriggersView) template() string {
	return "partials/triggers-view.html"
}

func (tsv *LabelTriggersView) data() interface{} {
	return nil
}

func (tsv *LabelTriggersView) scripts() []string {
	return nil
}

func (tsv *LabelTriggersView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tsv *LabelTriggersView) subRenderables() []Renderable {
	return tsv.subs
}

func (tsv *LabelTriggersView) renderPropagate() (*renderResult, error) {
	return renderPropagate(tsv)
}

func (tsv *LabelTriggersView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
