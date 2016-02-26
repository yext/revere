package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/web/vm"
)

type ProbeEdit struct {
	viewmodel *vm.Probe
	subs      map[string]Renderable
}

func NewProbeEdit(p *vm.Probe) *ProbeEdit {
	pe := ProbeEdit{}
	pe.viewmodel = p
	pe.subs = map[string]Renderable{}
	return &pe
}

func (pe *ProbeEdit) Template() string {
	tmpl, ok := pe.viewmodel.ProbeType().Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pe.viewmodel.Probe.ProbeType().Name()))
	}

	return path.Join(vm.ProbesDir, tmpl)
}

func (pe *ProbeEdit) Data() interface{} {
	return pe.viewmodel.Probe
}

func (pe *ProbeEdit) Scripts() []string {
	scripts := pe.viewmodel.ProbeType().Scripts()["edit"]

	return vm.AppendDir(vm.ProbesDir, scripts)
}

func (pe *ProbeEdit) Breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pe *ProbeEdit) SubRenderables() map[string]Renderable {
	return pe.subs
}

func (pe *ProbeEdit) RenderNow() bool {
	return true
}
