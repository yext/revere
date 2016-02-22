package vm

import (
	"fmt"
	"path"
)

type TargetEdit struct {
	viewmodel *Target
	subs      map[string]Renderable
}

func NewTargetEdit(t *Target) *TargetEdit {
	te := TargetEdit{}
	te.viewmodel = t
	te.subs = map[string]Renderable{}
	return &te
}

func (te *TargetEdit) Template() string {
	tmpl, ok := te.viewmodel.Templates()["edit"]
	if !ok {
		panic(fmt.Sprintf("Unable to find templates for target type %s", te.viewmodel.Target.TargetType().Name()))
	}

	return path.Join(targetsDir, tmpl)
}

func (te *TargetEdit) Data() interface{} {
	return te.viewmodel.Target
}

func (te *TargetEdit) Scripts() []string {
	scripts := te.viewmodel.Scripts()["edit"]

	return appendDir(targetsDir, scripts)
}

func (te *TargetEdit) Breadcrumbs() []Breadcrumb {
	return []Breadcrumb{}
}

func (te *TargetEdit) SubRenderables() map[string]Renderable {
	return te.subs
}

func (te *TargetEdit) RenderNow() bool {
	return true
}
