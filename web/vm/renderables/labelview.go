package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelView struct {
	label      *vm.Label
	subs       []Renderable
	saveStatus string
}

func NewLabelView(l *vm.Label, saveStatus []byte) *LabelView {
	lv := LabelView{}
	lv.label = l
	lv.subs = []Renderable{
		NewLabelTriggersView(l.Triggers),
		NewLabelMonitorsView(l.Monitors),
	}
	lv.saveStatus = string(saveStatus)
	return &lv
}

func (lv *LabelView) name() string {
	return "Label"
}

func (lv *LabelView) template() string {
	return "labels-view.html"
}

func (lv *LabelView) data() interface{} {
	return map[string]interface{}{
		"Label":      lv.label,
		"SaveStatus": lv.saveStatus,
	}
}

func (lv *LabelView) scripts() []string {
	return nil
}

func (lv *LabelView) breadcrumbs() []vm.Breadcrumb {
	return vm.LabelViewBcs(lv.label.Name, lv.label.Id())
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
