package vm

import (
	"html/template"
)

type Renderer interface {
	RenderView() (content template.HTML, scripts template.HTML, err error)
	RenderEdit() (content template.HTML, scripts template.HTML, err error)
}
