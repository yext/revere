package vm

import (
	"html/template"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
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
