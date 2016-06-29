package renderables

import "github.com/yext/revere/web/vm"

type MonitorLabelsEdit struct {
	monitorLabels []*vm.MonitorLabel
	allLabels     []*vm.Label
	subs          []Renderable
}

func NewMonitorLabelsEdit(mls []*vm.MonitorLabel, ls []*vm.Label) *MonitorLabelsEdit {
	mle := new(MonitorLabelsEdit)
	mle.monitorLabels = mls
	mle.allLabels = ls
	return mle
}

func (mle *MonitorLabelsEdit) name() string {
	return "MonitorLabels"
}

func (mle *MonitorLabelsEdit) template() string {
	return "partials/monitor-labels-edit.html"
}

func (mle *MonitorLabelsEdit) data() interface{} {
	return map[string]interface{}{
		"MonitorLabels": mle.monitorLabels,
		"AllLabels":     mle.allLabels,
	}
}

func (mle *MonitorLabelsEdit) scripts() []string {
	return []string{
		"component-list-edit.js",
		"monitor-labels-edit.js",
	}
}

func (mle *MonitorLabelsEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (mle *MonitorLabelsEdit) subRenderables() []Renderable {
	return mle.subs
}

func (mle *MonitorLabelsEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(mle)
}

func (mle *MonitorLabelsEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
