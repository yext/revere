package renderables

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

type Renderable interface {
	name() string
	template() string
	scripts() []string
	data() interface{}
	breadcrumbs() []vm.Breadcrumb
	subRenderables() []Renderable

	renderPropagate() (*renderResult, error)
	aggregatePipelineData(*renderResult, *renderResult)
}

type renderResult struct {
	name        string
	templates   []string
	scripts     []string
	data        map[string]interface{}
	breadcrumbs []vm.Breadcrumb
}

func Render(w io.Writer, r Renderable) error {
	result, err := r.renderPropagate()
	if err != nil {
		return err
	}

	return execute(w, result)
}

func RenderPartial(r Renderable) (template.HTML, error) {
	result, err := renderPropagateImmediate(r)
	if err != nil {
		return "", err
	}

	return template.HTML(fmt.Sprintf("%s", result.data["_Render"])), nil
}

func renderPropagate(r Renderable) (*renderResult, error) {
	parent := newRenderResult(r)

	for _, subrenderable := range r.subRenderables() {
		child, err := subrenderable.renderPropagate()
		if err != nil {
			return nil, err
		}

		parent.addSubRender(child)
		r.aggregatePipelineData(parent, child)
	}

	return parent, nil
}

func renderPropagateImmediate(r Renderable) (*renderResult, error) {
	result, err := renderPropagate(r)
	if err != nil {
		return nil, err
	}

	b := bytes.Buffer{}
	executeImmediate(&b, result)

	result.data["_Render"] = template.HTML(b.String())

	return result, nil
}

func newRenderResult(r Renderable) *renderResult {
	result := renderResult{}
	result.name = r.name()
	result.templates = []string{r.template()}
	result.scripts = r.scripts()
	result.data = map[string]interface{}{
		"_": r.data(),
	}
	result.breadcrumbs = r.breadcrumbs()
	return &result
}

func (current *renderResult) addSubRender(sub *renderResult) {
	current.templates = appendWithoutRepeat(current.templates, sub.templates...)
	current.scripts = appendWithoutRepeat(current.scripts, sub.scripts...)
	current.breadcrumbs = append(current.breadcrumbs, sub.breadcrumbs...)
}

func aggregatePipelineDataMap(parent *renderResult, child *renderResult) {
	parent.data[child.name] = child.data
}

func aggregatePipelineDataArray(parent *renderResult, child *renderResult) {
	if _, ok := parent.data["_Array"]; !ok {
		parent.data["_Array"] = []interface{}{}
	}
	array, ok := parent.data["_Array"].([]interface{})
	if !ok {
		panic("Non-array value found in \"_Array\" field of renderResult")
	}
	parent.data["_Array"] = append(array, child.data)
}

func appendWithoutRepeat(current []string, other ...string) []string {
	result := current
	tSet := map[string]bool{}

	for _, ele := range result {
		tSet[ele] = true
	}

	for _, ele := range other {
		if !tSet[ele] {
			tSet[ele] = true
			result = append(result, ele)
		}
	}

	return result
}

func prepareScripts(scripts []string) []string {
	length := len(scripts)
	result := make([]string, length, length)
	copy(result, scripts)

	for i, script := range scripts {
		result[i] = vm.GetScript(script)
	}

	return result
}

func prepareTemplates(templates []string) (*tmpl.Template, error) {
	if len(templates) == 0 {
		return nil, fmt.Errorf("Got error rendering views - no templates found")
	}

	t := tmpl.NewTemplate(templates[0])
	t.AddTmpls(templates[1:])

	return t, nil
}

func execute(w io.Writer, r *renderResult) error {
	t, err := prepareTemplates(r.templates)
	if err != nil {
		return err
	}

	data := r.data
	data["MoreScripts"] = prepareScripts(r.scripts)
	data["Breadcrumbs"] = r.breadcrumbs

	// Temporary to let footer / header tell between new / old templates
	data["New"] = "new"

	return t.Execute(w, data)
}

func executeImmediate(w io.Writer, r *renderResult) error {
	t, err := prepareTemplates(r.templates)
	if err != nil {
		return err
	}

	data := r.data["_"]

	return t.Execute(w, data)
}
