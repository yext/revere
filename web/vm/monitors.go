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
	view string = "monitors-view.html"
	edit string = "monitors-edit.html"
	dir  string = "web/views/"
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

func (vm *Monitor) render(tmplFile string) (template.HTML, error) {
	t := tmpl.NewTemplate(dir, tmplFile)
	b := bytes.Buffer{}
	err := t.Execute(&b, vm)
	if err != nil {
		return "", err
	}
	return template.HTML(b.String()), nil
}

func (vm *Monitor) RenderView() (template.HTML, error) {
	return vm.render(view)
}

func (vm *Monitor) RenderEdit() (template.HTML, error) {
	return vm.render(edit)
}
