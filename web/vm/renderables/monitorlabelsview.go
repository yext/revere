package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorLabelsView struct {
	monitorLabels []*vm.MonitorLabel
	subs          []Renderable
}

func NewMonitorLabelsView(mls []*vm.MonitorLabel) *MonitorLabelsView {
	mlv := new(MonitorLabelsView)
	mlv.monitorLabels = mls
	return mlv
}

func (mlv *MonitorLabelsView) name() string {
	return "MonitorLabels"
}

func (mlv *MonitorLabelsView) template() string {
	return "partials/monitor-labels-view.html"
}

func (mlv *MonitorLabelsView) data() interface{} {
	return map[string]interface{}{
		"MonitorLabels": mlv.monitorLabels,
	}
}

func (mlv *MonitorLabelsView) scripts() []string {
	return nil
}

func (mlv *MonitorLabelsView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (mlv *MonitorLabelsView) subRenderables() []Renderable {
	return mlv.subs
}

func (mlv *MonitorLabelsView) renderPropagate() (*renderResult, error) {
	return renderPropagate(mlv)
}

func (mlv *MonitorLabelsView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
