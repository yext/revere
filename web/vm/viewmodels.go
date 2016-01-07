package vm

import (
	"html/template"
)

type Renderer interface {
	RenderView() template.HTML
	RenderEdit() template.HTML
}
