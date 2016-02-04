package vm

import (
	"bytes"
	"html/template"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/targets"
	"github.com/yext/revere/web/tmpl"
)

type Monitor struct {
	Id            uint
	Name          string
	Owner         string
	Description   string
	Response      string
	Probe         probes.Probe
	ProbeTemplate template.HTML
	Changed       time.Time
	Version       int
	Archived      *time.Time
	Triggers      []*revere.Trigger
}

var (
	templateView string = "monitors-view.html"
	templateEdit string = "monitors-edit.html"

	scriptsView []string = []string{}
	scriptsEdit []string = []string{
		"revere.js",
		"monitors-edit.js",
		"probes/graphite-preview.js",
		"targets/email.js",
	}
)

func NewMonitor(m *revere.Monitor) (*Monitor, error) {
	vm := new(Monitor)

	vm.Id = m.Id
	vm.Name = m.Name
	vm.Owner = m.Owner
	vm.Description = m.Description
	vm.Response = m.Response
	vm.Changed = m.Changed
	vm.Version = m.Version
	vm.Archived = m.Archived
	vm.Triggers = m.Triggers

	// Load Probe
	probeType, err := probes.ProbeTypeById(m.ProbeType)
	if err != nil {
		return nil, err
	}

	vm.Probe, err = probeType.Load(m.ProbeJson)
	if err != nil {
		return nil, err
	}

	vm.ProbeTemplate, err = vm.Probe.Render()
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func BlankMonitor() (*Monitor, error) {
	vm := new(Monitor)

	vm.Triggers = []*revere.Trigger{
		&revere.Trigger{
			TargetTemplate: targets.DefaultTargetTemplate(),
		},
	}

	var err error
	vm.ProbeTemplate, err = probes.DefaultProbeTemplate()
	if err != nil {
		return nil, err
	}

	vm.Probe = probes.DefaultProbe()

	return vm, nil
}

func (vm *Monitor) render(templateFile string, scriptFiles []string) (content template.HTML, scripts template.HTML, err error) {
	t := tmpl.NewTemplate(templateFile)
	b := bytes.Buffer{}
	err = t.Execute(&b, vm)
	if err != nil {
		return "", "", err
	}
	content = template.HTML(b.String())
	scripts = newScripts(scriptFiles)
	return content, scripts, nil
}

func (vm *Monitor) RenderView() (content template.HTML, scripts template.HTML, err error) {
	return vm.render(templateView, scriptsView)
}

func (vm *Monitor) RenderEdit() (content template.HTML, scripts template.HTML, err error) {
	return vm.render(templateEdit, scriptsEdit)
}
