package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorTriggerEdit struct {
	viewmodel *vm.MonitorTrigger
	subs      []Renderable
}

func NewMonitorTriggerEdit(mt *vm.MonitorTrigger) *MonitorTriggerEdit {
	te := MonitorTriggerEdit{}
	te.viewmodel = mt
	te.subs = []Renderable{
		NewTargetEdit(mt.Trigger.Target),
	}
	return &te
}

func (te *MonitorTriggerEdit) name() string {
	return "Trigger"
}

func (te *MonitorTriggerEdit) template() string {
	return "partials/trigger-edit.html"
}

func (te *MonitorTriggerEdit) data() interface{} {
	return nil
}

func (te *MonitorTriggerEdit) scripts() []string {
	return nil
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
