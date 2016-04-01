package renderables

import (
	"github.com/yext/revere/web/vm"
)

type LabelsIndex struct {
	labels []*vm.Label
	subs   []Renderable
}

func NewLabelsIndex(ss []*vm.Label) *LabelsIndex {
	lsi := new(LabelsIndex)
	lsi.labels = ss

	return lsi
}

func (lsi *LabelsIndex) name() string {
	return "LabelsIndex"
}

func (lsi *LabelsIndex) template() string {
	return "labels-index.html"
}

func (lsi *LabelsIndex) data() interface{} {
	return lsi.labels
}

func (lsi *LabelsIndex) scripts() []string {
	return nil
}

func (lsi *LabelsIndex) breadcrumbs() []vm.Breadcrumb {
	return vm.LabelIndexBcs()
}

func (lsi *LabelsIndex) subRenderables() []Renderable {
	return nil
}

func (lsi *LabelsIndex) renderPropagate() (*renderResult, error) {
	return renderPropagate(lsi)
}

func (lsi *LabelsIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
