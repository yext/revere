package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorTriggersView struct {
	viewmodel []*vm.MonitorTrigger
	subs      []Renderable
}

func NewMonitorTriggersView(mts []*vm.MonitorTrigger) *MonitorTriggersView {
	tsv := MonitorTriggersView{}
	for _, monitortrigger := range mts {
		tsv.subs = append(tsv.subs, NewMonitorTriggerView(monitortrigger))
	}

	return &tsv
}

func (tsv *MonitorTriggersView) name() string {
	return "Triggers"
}

func (tsv *MonitorTriggersView) template() string {
	return "partials/triggers-view.html"
}

func (tsv *MonitorTriggersView) data() interface{} {
	return tsv.viewmodel
}

func (tsv *MonitorTriggersView) scripts() []string {
	return nil
}

func (tsv *MonitorTriggersView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (tsv *MonitorTriggersView) subRenderables() []Renderable {
	return tsv.subs
}

func (tsv *MonitorTriggersView) renderPropagate() (*renderResult, error) {
	return renderPropagate(tsv)
}

func (tsv *MonitorTriggersView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
