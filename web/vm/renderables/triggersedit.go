package renderables

import (
	"github.com/yext/revere/web/vm"
)

type TriggersEdit struct {
	viewmodel []*vm.Trigger
	subs      []Renderable
}

func NewTriggersEdit(ts []*vm.Trigger) *TriggersEdit {
	tse := TriggersEdit{}
	for _, trigger := range ts {
		tse.subs = append(tse.subs, NewTriggerEdit(trigger))
	}

	return &tse
}

func (tse *TriggersEdit) name() string {
	return "Triggers"
}

func (tse *TriggersEdit) template() string {
	return "partials/triggers-edit.html"
}

func (tse *TriggersEdit) data() interface{} {
	return tse.viewmodel
}

func (tse *TriggersEdit) scripts() []string {
	return []string{
		"triggers-edit.js",
	}
}

func (tse *TriggersEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tse *TriggersEdit) subRenderables() []Renderable {
	return tse.subs
}

func (tse *TriggersEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(tse)
}

func (tse *TriggersEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
