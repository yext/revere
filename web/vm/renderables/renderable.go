/*
Package renderables processes package vm structs and bundles them with the
appropriate metadata, html templates, and scripts to create a fully functional
web page.

Generally, Revere pages are created out of a top-level "viewmodel" Renderable
interface that provides "sub-Renderables" out of the components it needs to
display.

Rendering is done recursively from the top-level Renderable down to all of its
sub-Renderables.

First, all sub-Renderables are recursively processed to create respective
renderResults, which include the names of those Renderables (that the templates
will use to access related data), the template files needed to build the
skeleton HTML, javascript files, nested data objects to plug into the
templates, and breadcrumbs.

Next, the sub-Renderables' renderResults are combined with those of the parent
Renderable. This is done by first simply appending the child's templates,
scripts, and breadcrumbs to those of the parent.  The implementer specifies how
the data will be added to the parent data - most Revere Renderables will use
aggregatePipelineDataMap, which will add the child renderResult data as map
entries into the parent renderResult, using child.name as the key. In case the
template needs to access data in an array, sub-Renderable data can be
aggregated using aggregatePipelineDataArray, which places all of the child data
into a interface{} array in order of the corresponding positions in the
sub-Renderables array, and added to the parent data object under the key
"_Array".

Lastly the renderResult is returned (if the regular renderPropogate method is
specified, which should be the default in most cases). In the cases of Probes,
Resources, Targets, and Settings, go's template package will not allow
dynamically-specified template files and requires an alternate solution.
Calling renderPropagateImmediate on the renderPropagate step renders the entire
Renderable into a template.HTML string, which is then added to the data object
of the renderResult under the key "_Render".  This way, a template can be
dynamically specified, rendered into HTML, and inserted into the parent HTML
template.
*/
package renderables

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/juju/errors"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
)

// The Renderable interface allows implementations to specify all required
// files in order to generate a web page.
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

// Render constructs an HTML document using a Renderable and writes it to the
// specified io.Writer.
func Render(w io.Writer, r Renderable) error {
	result, err := r.renderPropagate()
	if err != nil {
		return errors.Trace(err)
	}

	return execute(w, result)
}

// RenderPartial constructs an HTML document using a Renderable, but returns
// the HTML as a string instead of writing it to an io.Writer.
func RenderPartial(r Renderable) (template.HTML, error) {
	result, err := renderPropagateImmediate(r)
	if err != nil {
		return "", errors.Trace(err)
	}

	return template.HTML(fmt.Sprintf("%s", result.data["_Render"])), nil
}

func renderPropagate(r Renderable) (*renderResult, error) {
	parent := newRenderResult(r)

	for _, subrenderable := range r.subRenderables() {
		child, err := subrenderable.renderPropagate()
		if err != nil {
			return nil, errors.Trace(err)
		}

		parent.addSubRender(child)
		r.aggregatePipelineData(parent, child)
	}

	return parent, nil
}

func renderPropagateImmediate(r Renderable) (*renderResult, error) {
	result, err := renderPropagate(r)
	if err != nil {
		return nil, errors.Trace(err)
	}

	b := bytes.Buffer{}
	err = executeImmediate(&b, result)
	if err != nil {
		return nil, errors.Trace(err)
	}

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
		result[i] = tmpl.GetScript(script)
	}

	return result
}

func prepareTemplates(templates []string) (*tmpl.Template, error) {
	if len(templates) == 0 {
		return nil, errors.Errorf("Got error rendering views - no templates found")
	}

	t := tmpl.NewTemplate(templates[0])
	t.AddTmpls(templates[1:])

	return t, nil
}

func execute(w io.Writer, r *renderResult) error {
	t, err := prepareTemplates(r.templates)
	if err != nil {
		return errors.Trace(err)
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
		return errors.Trace(err)
	}

	data := r.data["_"]

	return t.Execute(w, data)
}
