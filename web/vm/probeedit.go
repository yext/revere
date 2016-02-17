package vm

import (
	"fmt"
	"path"
)

type ProbeEdit struct {
	viewmodel *Probe
	subs      map[string]Renderable
}

func NewProbeEdit(p *Probe) *ProbeEdit {
	pe := ProbeEdit{}
	pe.viewmodel = p
	pe.subs = map[string]Renderable{}
	return &pe
}

func (pe *ProbeEdit) Template() string {
	tmpl, ok := pe.viewmodel.Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for probe type %s", pe.viewmodel.Probe.ProbeType().Name()))
	}

	return path.Join(probesDir, tmpl)
}

func (pe *ProbeEdit) Data() interface{} {
	return pe.viewmodel.Probe
}

func (pe *ProbeEdit) Scripts() []string {
	scripts, ok := pe.viewmodel.Scripts()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find scripts for probe type %s", pe.viewmodel.Probe.ProbeType().Name()))
	}

	return appendDir(probesDir, scripts)
}

func (pe *ProbeEdit) Breadcrumbs() []Breadcrumb {
	return []Breadcrumb{}
}

func (pe *ProbeEdit) SubRenderables() map[string]Renderable {
	return pe.subs
}

func (pe *ProbeEdit) RenderNow() bool {
	return true
}
