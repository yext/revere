package renderables

import (
	"fmt"
	"path"

	"github.com/yext/revere/web/vm"
)

type ProbeView struct {
	viewmodel *vm.Probe
	subs      map[string]Renderable
}

func NewProbeView(p *vm.Probe) *ProbeView {
	pv := ProbeView{}
	pv.viewmodel = p
	pv.subs = map[string]Renderable{}
	return &pv
}

func (pv *ProbeView) Template() string {
	tmpl, ok := pv.viewmodel.ProbeType().Templates()["view"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pv.viewmodel.Probe.ProbeType().Name()))
	}

	return path.Join(vm.ProbesDir, tmpl)
}

func (pv *ProbeView) Data() interface{} {
	return pv.viewmodel.Probe
}

func (pv *ProbeView) Scripts() []string {
	scripts := pv.viewmodel.ProbeType().Scripts()["view"]

	return vm.AppendDir(vm.ProbesDir, scripts)
}

func (pv *ProbeView) Breadcrumbs() []vm.Breadcrumb {
	return []vm.Breadcrumb{}
}

func (pv *ProbeView) SubRenderables() map[string]Renderable {
	return pv.subs
}

func (pv *ProbeView) RenderNow() bool {
	return true
}
