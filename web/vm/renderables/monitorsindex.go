package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorsIndex struct {
	monitors []*vm.Monitor
	labels   []*vm.Label
	subs     []Renderable
}

func NewMonitorsIndex(ms []*vm.Monitor, ls []*vm.Label) *MonitorsIndex {
	msi := new(MonitorsIndex)
	msi.monitors = ms
	msi.labels = ls

	return msi
}

func (msi *MonitorsIndex) name() string {
	return "MonitorsIndex"
}

func (msi *MonitorsIndex) template() string {
	return "monitors-index.html"
}

func (msi *MonitorsIndex) data() interface{} {
	return map[string]interface{}{
		"Monitors": msi.monitors,
		"Labels":   msi.labels,
	}
}

func (msi *MonitorsIndex) scripts() []string {
	return []string{
		"label-filter.js",
		"monitors-index.js",
	}
}

func (msi *MonitorsIndex) breadcrumbs() []vm.Breadcrumb {
	return vm.MonitorIndexBcs()
}

func (msi *MonitorsIndex) subRenderables() []Renderable {
	return nil
}

func (msi *MonitorsIndex) renderPropagate() (*renderResult, error) {
	return renderPropagate(msi)
}

func (msi *MonitorsIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
