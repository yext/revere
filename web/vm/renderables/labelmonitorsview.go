package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelMonitorsView struct {
	viewmodel interface{}
	subs      []Renderable
}

func NewLabelMonitorsView(lms interface{}) *LabelMonitorsView {
	lmv := new(LabelMonitorsView)
	lmv.viewmodel = lms
	return lmv
}

func (lmv *LabelMonitorsView) name() string {
	return "LabelMonitors"
}

func (lmv *LabelMonitorsView) template() string {
	return "partials/label-monitors-view.html"
}

func (lmv *LabelMonitorsView) data() interface{} {
	return lmv.viewmodel
}

func (lmv *LabelMonitorsView) scripts() []string {
	return nil
}

func (lmv *LabelMonitorsView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (lmv *LabelMonitorsView) subRenderables() []Renderable {
	return lmv.subs
}

func (lmv *LabelMonitorsView) renderPropagate() (*renderResult, error) {
	return renderPropagate(lmv)
}

func (lmv *LabelMonitorsView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
