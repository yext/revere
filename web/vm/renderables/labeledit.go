package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelEdit struct {
	viewmodel *vm.Label
	subs      []Renderable
}

func NewLabelEdit(l *vm.Label) *LabelEdit {
	le := LabelEdit{}
	le.viewmodel = l
	le.subs = []Renderable{
	// "Triggers": le.TriggersEdit(m.Triggers),
	}
	return &le
}

func (mv *LabelEdit) name() string {
	return "Label"
}

func (le *LabelEdit) template() string {
	return "labels-edit.html"
}

func (le *LabelEdit) data() interface{} {
	return le.viewmodel
}

func (le *LabelEdit) scripts() []string {
	return []string{
		"revere.js",
		"labels-edit.js",
	}
}

func (le *LabelEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (le *LabelEdit) subRenderables() []Renderable {
	return le.subs
}

func (le *LabelEdit) renderPropogate() (*renderResult, error) {
	return renderPropogate(le)
}

func (le *LabelEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
