package renderables

import (
	"path"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/web/vm"
)

type DataSourceView struct {
	datasource *datasources.VM
}

func NewDataSourceView(ds *datasources.VM) *DataSourceView {
	dsv := DataSourceView{}
	dsv.datasource = ds
	return &dsv
}

func (dsv *DataSourceView) name() string {
	return dsv.datasource.Name()
}

func (dsv *DataSourceView) template() string {
	return path.Join(datasources.DataSourceDir, dsv.datasource.Templates())
}

func (dsv *DataSourceView) data() interface{} {
	return dsv.datasource
}

func (dsv *DataSourceView) scripts() []string {
	return vm.AppendDir(datasources.DataSourceDir, dsv.datasource.Scripts())
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
