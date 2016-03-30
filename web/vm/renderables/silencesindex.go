package renderables

import (
	"github.com/yext/revere/web/vm"
)

type SilencesIndex struct {
	silences []*vm.Silence
	subs     []Renderable
}

func NewSilencesIndex(ss []*vm.Silence) *SilencesIndex {
	ssi := new(SilencesIndex)
	ssi.silences = ss

	return ssi
}

func (ssi *SilencesIndex) name() string {
	return "SilencesIndex"
}

func (ssi *SilencesIndex) template() string {
	return "silences-index.html"
}

func (ssi *SilencesIndex) data() interface{} {
	return ssi.silences
}

func (ssi *SilencesIndex) scripts() []string {
	return []string{
		"silences-index.js",
	}
}

func (ssi *SilencesIndex) breadcrumbs() []vm.Breadcrumb {
	return vm.SilencesIndexBcs()
}

func (ssi *SilencesIndex) subRenderables() []Renderable {
	return nil
}

func (ssi *SilencesIndex) renderPropagate() (*renderResult, error) {
	return renderPropagate(ssi)
}

func (ssi *SilencesIndex) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
