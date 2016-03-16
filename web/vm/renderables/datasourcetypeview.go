package renderables

import (
	"path"

	"github.com/yext/revere/web/vm"
)

type DataSourceTypeView struct {
	*vm.DataSourceTypeViewModel
}

func NewDataSourceTypeView(dst *vm.DataSourceTypeViewModel) *DataSourceTypeView {
	dstv := DataSourceTypeView{}
	dstv.DataSourceTypeViewModel = dst
	return &dstv
}

func (dst *DataSourceTypeView) name() string {
	return dst.Type.Name()
}

func (dst *DataSourceTypeView) template() string {
	return path.Join(vm.DataSourceDir, dst.Type.Template())
}

func (dst *DataSourceTypeView) data() interface{} {
	return dst.DataSourceTypeViewModel
}

func (dst *DataSourceTypeView) scripts() []string {
	return vm.AppendDir(vm.DataSourceDir, dst.Type.Scripts())
}

func (dst *DataSourceTypeView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (dst *DataSourceTypeView) subRenderables() []Renderable {
	return nil
}

func (dst *DataSourceTypeView) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(dst)
}

func (dst *DataSourceTypeView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
