package vm

type Renderable interface {
	Template() string
	Scripts() []string
	Data() interface{}
	Breadcrumbs() []Breadcrumb
	SubRenderables() map[string]Renderable
}

type RenderResult struct {
	Templates   []string
	Scripts     []string
	Data        map[string]interface{}
	Breadcrumbs []Breadcrumb
}

func (current *RenderResult) AddSubRender(name string, sub *RenderResult) {
	current.Templates = append(current.Templates, sub.Templates...)
	current.Scripts = append(current.Scripts, sub.Scripts...)
	current.Data[name] = sub.Data["Data"]
	current.Breadcrumbs = append(current.Breadcrumbs, sub.Breadcrumbs...)
}

func NewRenderResult(r Renderable) *RenderResult {
	result := RenderResult{}
	result.Templates = []string{r.Template()}
	result.Scripts = r.Scripts()
	result.Data = map[string]interface{}{
		"Data": r.Data(),
	}
	result.Breadcrumbs = []Breadcrumb{}
	return &result
}
