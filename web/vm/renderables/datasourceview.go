package renderables

import (
	"path"

	"github.com/yext/revere/datasource"
	"github.com/yext/revere/web/vm"
)

type DataSourceView struct {
	datasource *datasource.VM
}

func NewDataSourceView(ds *datasource.VM) *DataSourceView {
	dsv := DataSourceView{}
	dsv.datasource = ds
	return &dsv
}

func (dsv *DataSourceView) name() string {
	return dsv.datasource.Name()
}

func (dsv *DataSourceView) template() string {
	return path.Join(datasource.DataSourceDir, dsv.datasource.Templates())
}

func (dsv *DataSourceView) data() interface{} {
	return dsv.datasource
}

func (dsv *DataSourceView) scripts() []string {
	return nil
}

func (dsv *DataSourceView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (dsv *DataSourceView) subRenderables() []Renderable {
	return nil
}

func (dsv *DataSourceView) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(dsv)
}

func (dsv *DataSourceView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
