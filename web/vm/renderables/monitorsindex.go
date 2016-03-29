package renderables

import (
	"github.com/yext/revere/web/vm"
)

type MonitorsIndex struct {
	viewmodel []*vm.Monitor
	subs      []Renderable
}

func NewMonitorsIndex(ms []*vm.Monitor) *MonitorsIndex {
	msi := new(MonitorsIndex)
	msi.viewmodel = ms

	return msi
}

func (msi *MonitorsIndex) name() string {
	return "MonitorsIndex"
}

func (msi *MonitorsIndex) template() string {
	return "monitors-index.html"
}

func (msi *MonitorsIndex) data() interface{} {
	return msi.viewmodel
}

func (msi *MonitorsIndex) scripts() []string {
	return nil
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
