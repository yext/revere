package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorTriggersEdit struct {
	subs []Renderable
}

func NewMonitorTriggersEdit(mts []*vm.MonitorTrigger) *MonitorTriggersEdit {
	tse := MonitorTriggersEdit{}
	tse.subs = []Renderable{
		NewMonitorTriggerEdit(vm.BlankMonitorTrigger()),
	}
	for _, monitortrigger := range mts {
		tse.subs = append(tse.subs, NewMonitorTriggerEdit(monitortrigger))
	}

	return &tse
}

func (tse *MonitorTriggersEdit) name() string {
	return "Triggers"
}

func (tse *MonitorTriggersEdit) template() string {
	return "partials/monitor-triggers-edit.html"
}

func (tse *MonitorTriggersEdit) data() interface{} {
	return nil
}

func (tse *MonitorTriggersEdit) scripts() []string {
	return []string{
		"monitor-triggers-edit.js",
		"triggers-edit.js",
	}
}

func (tse *MonitorTriggersEdit) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tse *MonitorTriggersEdit) subRenderables() []Renderable {
	return tse.subs
}

func (tse *MonitorTriggersEdit) renderPropagate() (*renderResult, error) {
	return renderPropagate(tse)
}

func (tse *MonitorTriggersEdit) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
