package renderables

import (
	"github.com/yext/revere/db"
	"github.com/yext/revere/web/vm"
)

type ActiveIssues struct {
	labels        []*vm.Label
	subprobes     []*vm.Subprobe
	monitorLabels map[db.MonitorID][]*vm.MonitorLabel
	subs          []Renderable
}

func NewActiveIssues(ss []*vm.Subprobe, ls []*vm.Label, mls map[db.MonitorID][]*vm.MonitorLabel) *ActiveIssues {
	return &ActiveIssues{ls, ss, mls, nil}
}

func (ai *ActiveIssues) name() string {
	return "ActiveIssues"
}

func (ai *ActiveIssues) template() string {
	return "active-issues.html"
}

func (ai *ActiveIssues) data() interface{} {
	return map[string]interface{}{
		"Labels":        ai.labels,
		"Subprobes":     ai.subprobes,
		"MonitorLabels": ai.monitorLabels,
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
	return nil
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
