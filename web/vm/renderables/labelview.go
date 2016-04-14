package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelView struct {
	viewmodel *vm.Label
	subs      []Renderable
}

func NewLabelView(l *vm.Label) *LabelView {
	lv := LabelView{}
	lv.viewmodel = l
	lv.subs = []Renderable{
		NewTriggersView(l.Triggers),
		NewLabelMonitorsView(l.Monitors),
	}
	return &lv
}

func (lv *LabelView) name() string {
	return "Label"
}

func (lv *LabelView) template() string {
	return "labels-view.html"
}

func (lv *LabelView) data() interface{} {
	return lv.viewmodel
}

func (lv *LabelView) scripts() []string {
	return nil
}

func (lv *LabelView) breadcrumbs() []vm.Breadcrumb {
	return vm.LabelViewBcs(lv.viewmodel.Name(), lv.viewmodel.Id())
}

func (lv *LabelView) subRenderables() []Renderable {
	return lv.subs
}

func (lv *LabelView) renderPropagate() (*renderResult, error) {
	return renderPropagate(lv)
}

func (lv *LabelView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
