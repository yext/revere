package renderables

import "github.com/yext/revere/web/vm"

type LabelMonitorsEdit struct {
	labelMonitors []*vm.LabelMonitor
	allMonitors   []*vm.Monitor
	subs          []Renderable
}

func NewLabelMonitorsEdit(lms []*vm.LabelMonitor, ms []*vm.Monitor) *LabelMonitorsEdit {
	lme := new(LabelMonitorsEdit)
	lme.labelMonitors = lms
	lme.allMonitors = ms
	return lme
}

func (lme *LabelMonitorsEdit) name() string {
	return "LabelMonitors"
}

func (lme *LabelMonitorsEdit) template() string {
	return "partials/label-monitors-edit.html"
}

func (lme *LabelMonitorsEdit) data() interface{} {
	return map[string]interface{}{
		"LabelMonitors": lme.labelMonitors,
		"AllMonitors":   lme.allMonitors,
	}
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
