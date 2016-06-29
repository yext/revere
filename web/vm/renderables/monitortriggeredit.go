package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorTriggerEdit struct {
	monitorTrigger *vm.MonitorTrigger
	subs           []Renderable
}

func NewMonitorTriggerEdit(t *vm.MonitorTrigger) *MonitorTriggerEdit {
	te := MonitorTriggerEdit{}
	te.monitorTrigger = t
	te.subs = []Renderable{
		NewTargetEdit(t.Trigger.Target),
	}
	return &te
}

func (te *MonitorTriggerEdit) name() string {
	return "Trigger"
}

func (te *MonitorTriggerEdit) template() string {
	return "partials/monitor-trigger-edit.html"
}

func (te *MonitorTriggerEdit) data() interface{} {
	return te.monitorTrigger
}

func (te *MonitorTriggerEdit) scripts() []string {
	return []string{
		"trigger-edit.js",
	}
}

func (te *MonitorTriggerEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (te *MonitorTriggerEdit) subRenderables() []Renderable {
	return te.subs
}

func (te *MonitorTriggerEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(te)
}

func (te *MonitorTriggerEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataMap(parent, child)
}
