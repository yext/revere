package vm

import (
	"html/template"
)

type Renderer interface {
	RenderView() (template.HTML, error)
	RenderEdit() (template.HTML, error)
}
