package renderables

import (
	"html/template"

	"github.com/yext/revere/web/vm"
)

type Renderable interface {
	Template() string
	Scripts() []string
	Data() interface{}
	Breadcrumbs() []vm.Breadcrumb
	SubRenderables() map[string]Renderable
	RenderNow() bool
}

type RenderResult struct {
	Templates   []string
	Scripts     []string
	Data        map[string]interface{}
	Breadcrumbs []vm.Breadcrumb
}

func (current *RenderResult) AddSubRender(name string, sub *RenderResult) {
	current.Templates = append(current.Templates, sub.Templates...)
	current.Scripts = append(current.Scripts, sub.Scripts...)
	current.Data[name] = sub.Data
	current.Breadcrumbs = append(current.Breadcrumbs, sub.Breadcrumbs...)
}

func (current *RenderResult) AddRendered(name string, renderedHtml template.HTML) {
	current.Data[name] = renderedHtml
}

func NewRenderResult(r Renderable) *RenderResult {
	result := RenderResult{}
	result.Templates = []string{r.Template()}
	result.Scripts = r.Scripts()
	result.Data = map[string]interface{}{
		"Data": r.Data(),
	}
	result.Breadcrumbs = []vm.Breadcrumb{}
	return &result
}