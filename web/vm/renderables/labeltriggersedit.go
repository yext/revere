package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelTriggersEdit struct {
	viewmodel []*vm.Trigger
	subs      []Renderable
}

func NewLabelTriggersEdit(ts []*vm.LabelTrigger) *LabelTriggersEdit {
	tse := LabelTriggersEdit{}
	tse.subs = []Renderable{
		NewLabelTriggerEdit(vm.BlankLabelTrigger()),
	}
	for _, labeltrigger := range ts {
		tse.subs = append(tse.subs, NewLabelTriggerEdit(labeltrigger))
	}

	return &tse
}

func (tse *LabelTriggersEdit) name() string {
	return "Triggers"
}

func (tse *LabelTriggersEdit) template() string {
	return "partials/triggers-edit.html"
}

func (tse *LabelTriggersEdit) data() interface{} {
	return tse.viewmodel
}

func (tse *LabelTriggersEdit) scripts() []string {
	return []string{
		"triggers-edit.js",
	}
}

func (tse *LabelTriggersEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tse *LabelTriggersEdit) subRenderables() []Renderable {
	return tse.subs
}

func (tse *LabelTriggersEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(tse)
}

func (tse *LabelTriggersEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
