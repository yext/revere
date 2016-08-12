package renderables

import (
	"path"

	"github.com/yext/revere/resource"
	"github.com/yext/revere/web/vm"
)

type ResourceView struct {
	resource *resource.VM
}

func NewResourceView(ds *resource.VM) *ResourceView {
	dsv := ResourceView{}
	dsv.resource = ds
	return &dsv
}

func (dsv *ResourceView) name() string {
	return dsv.resource.Name()
}

func (dsv *ResourceView) template() string {
	return path.Join(resource.ResourceDir, dsv.resource.Templates())
}

func (dsv *ResourceView) data() interface{} {
	return dsv.resource
}

func (dsv *ResourceView) scripts() []string {
	return nil
}

func (dsv *ResourceView) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (dsv *ResourceView) subRenderables() []Renderable {
	return nil
}

func (dsv *ResourceView) renderPropagate() (*renderResult, error) {
	return renderPropagateImmediate(dsv)
}

func (dsv *ResourceView) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
