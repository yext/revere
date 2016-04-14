package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelEdit struct {
	viewmodel *vm.Label
	subs      []Renderable
}

func NewLabelEdit(l *vm.Label, ms []*vm.Monitor) *LabelEdit {
	le := LabelEdit{}
	le.viewmodel = l
	le.subs = []Renderable{
		NewTriggersEdit(l.Triggers),
		NewLabelMonitorsEdit(l.Monitors, ms),
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
		"labels-edit.js",
	}
}

func (le *LabelEdit) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (le *LabelEdit) subRenderables() []Renderable {
	return le.subs
}

func (le *LabelEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(le)
}

func (le *LabelEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
