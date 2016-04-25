package renderables

import (
	"github.com/yext/revere/web/vm"
)

type DataSourceIndex struct {
	viewmodels []*vm.DataSourceType
	subs       []Renderable
}

func NewDataSourceIndex(dstvms []*vm.DataSourceType) *DataSourceIndex {
	dsi := DataSourceIndex{}
	dsi.viewmodels = dstvms
	dsi.subs = make([]Renderable, 0)
	for _, dstvm := range dstvms {
		dsi.subs = append(dsi.subs, NewDataSourceTypeView(dstvm))
	}
	return &dsi
}

func (dsi *DataSourceIndex) name() string {
	return "Data Sources"
}

func (dsi *DataSourceIndex) template() string {
	return "datasources-index.html"
}

func (dsi *DataSourceIndex) data() interface{} {
	return dsi.viewmodels
}

func (dsi *DataSourceIndex) scripts() []string {
	return []string{
		"datasources.js",
	}
}

func (dsi *DataSourceIndex) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (dsi *DataSourceIndex) subRenderables() []Renderable {
	return dsi.subs
}

func (dsi *DataSourceIndex) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(dsi)
}

func (dsi *DataSourceIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
