package renderables

import (
	"github.com/yext/revere/resource"
	"github.com/yext/revere/web/vm"
)

type ResourcesIndex struct {
	resources  []*resource.VM
	subs       []Renderable
	saveStatus string
}

func NewResourcesIndex(dss []*resource.VM, saveStatus []byte) *ResourcesIndex {
	ri := ResourcesIndex{}
	ri.resources = dss
	ri.subs = make([]Renderable, len(dss))
	for i, ds := range dss {
		ri.subs[i] = NewResourceView(ds)
	}
	ri.saveStatus = string(saveStatus)
	return &ri
}

func (ri *ResourcesIndex) name() string {
	return "Resources"
}

func (ri *ResourcesIndex) template() string {
	return "resources-index.html"
}

func (ri *ResourcesIndex) data() interface{} {
	return map[string]interface{}{
		"Resources":  resource.AllTypes(),
		"SaveStatus": ri.saveStatus,
	}
}

func (ri *ResourcesIndex) scripts() []string {
	return resource.AllScripts()
}

func (ri *ResourcesIndex) breadcrumbs() []vm.Breadcrumb {
	return nil
}

func (ri *ResourcesIndex) subRenderables() []Renderable {
	return ri.subs
}

func (ri *ResourcesIndex) renderPropagate() (*renderResult, error) {
	return renderPropagate(ri)
}

func (ri *ResourcesIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
