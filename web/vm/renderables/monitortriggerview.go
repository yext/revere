package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorTriggerView struct {
	viewmodel *vm.MonitorTrigger
	subs      []Renderable
}

func NewMonitorTriggerView(mt *vm.MonitorTrigger) *MonitorTriggerView {
	tv := MonitorTriggerView{}
	tv.viewmodel = mt
	tv.subs = []Renderable{
		NewTargetView(mt.Trigger.Target),
	}
	return &tv
}

func (tv *MonitorTriggerView) name() string {
	return "Trigger"
}

func (tv *MonitorTriggerView) template() string {
	return "partials/trigger-view.html"
}

func (tv *MonitorTriggerView) data() interface{} {
	return tv.viewmodel
}

func (tv *MonitorTriggerView) scripts() []string {
	return nil
}

func (tv *MonitorTriggerView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tv *MonitorTriggerView) subRenderables() []Renderable {
	return tv.subs
}

func (tv *MonitorTriggerView) renderPropagate() (*renderResult, error) {
	return renderPropagate(tv)
}

func (tv *MonitorTriggerView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
