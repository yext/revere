package vm

import (
	"fmt"
	"path"
)

type ProbeView struct {
	viewmodel *Probe
	subs      map[string]Renderable
}

func NewProbeView(p *Probe) *ProbeView {
	pv := ProbeView{}
	pv.viewmodel = p
	pv.subs = map[string]Renderable{}
	return &pv
}

func (pv *ProbeView) Template() string {
	tmpl, ok := pv.viewmodel.Templates()["view"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pv.viewmodel.Probe.ProbeType().Name()))
	}

	return path.Join(probesDir, tmpl)
}

func (pv *ProbeView) Data() interface{} {
	return pv.viewmodel.Probe
}

func (pv *ProbeView) Scripts() []string {
	scripts, ok := pv.viewmodel.Scripts()["view"]
	if !ok {
		panic(fmt.Sprintf("Unable to find scripts for probe type %s", pv.viewmodel.Probe.ProbeType().Name()))
	}

	return appendDir(probesDir, scripts)
}

func (pv *ProbeView) Breadcrumbs() []Breadcrumb {
	return []Breadcrumb{}
}

func (pv *ProbeView) SubRenderables() map[string]Renderable {
	return pv.subs
}

func (pv *ProbeView) RenderNow() bool {
	return true
}
