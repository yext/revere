package renderables

import (
	"github.com/yext/revere"
	"github.com/yext/revere/web/vm"
)

type ActiveIssues struct {
	labels    []*vm.Label
	subprobes []*vm.Subprobe
	subs      []Renderable
}

func NewActiveIssues(ss []*vm.Subprobe, ls []*vm.Label) *ActiveIssues {
	ai := new(ActiveIssues)
	ai.labels = ls
	ai.subprobes = ss

	return ai
}

func (ai *ActiveIssues) name() string {
	return "ActiveIssues"
}

func (ai *ActiveIssues) template() string {
	return "active-issues.html"
}

func (ai *ActiveIssues) data() interface{} {
	return map[string]interface{}{
		"Labels":    ai.labels,
		"Subprobes": ai.subprobes,
		"States":    revere.ReverseStates,
	}
}

func (ai *ActiveIssues) scripts() []string {
	return []string{
		"label-filter.js",
		"entered-states.js",
		"active-issues.js",
	}
}

func (ai *ActiveIssues) breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{vm.Breadcrumb{"active issues", "/"}}
}

func (ai *ActiveIssues) subRenderables() []Renderable {
	return nil
}

func (ai *ActiveIssues) renderPropagate() (*renderResult, error) {
	return renderPropagate(ai)
}

func (ai *ActiveIssues) aggregatePipelineData(parent *renderResult, child *renderResult) {
	aggregatePipelineDataArray(parent, child)
}
