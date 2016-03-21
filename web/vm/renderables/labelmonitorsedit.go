package renderables

import "github.com/yext/revere/web/vm"

type LabelMonitorsEdit struct {
	viewmodel interface{}
	subs      []Renderable
}

func NewLabelMonitorsEdit(lms interface{}) *LabelMonitorsEdit {
	lme := new(LabelMonitorsEdit)
	lme.viewmodel = lms
	return lme
}

func (lme *LabelMonitorsEdit) name() string {
	return "LabelMonitors"
}

func (lme *LabelMonitorsEdit) template() string {
	return "partials/label-monitors-edit.html"
}

func (lme *LabelMonitorsEdit) data() interface{} {
	return lme.viewmodel
}

func (lme *LabelMonitorsEdit) scripts() []string {
	return []string{
		"label-monitors-edit.js",
	}
}

func (lme *LabelMonitorsEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (lme *LabelMonitorsEdit) subRenderables() []Renderable {
	return lme.subs
}

func (lme *LabelMonitorsEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(lme)
}

func (lme *LabelMonitorsEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
