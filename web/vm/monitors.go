package vm

import (
	"html/template"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/targets"
)

type Monitor struct {
	Id            uint
	Name          string
	Owner         string
	Description   string
	Response      string
	Probe         *Probe
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
	viewmodel := new(Monitor)

	viewmodel.Id = m.Id
	viewmodel.Name = m.Name
	viewmodel.Owner = m.Owner
	viewmodel.Description = m.Description
	viewmodel.Response = m.Response
	viewmodel.Changed = m.Changed
	viewmodel.Version = m.Version
	viewmodel.Archived = m.Archived
	viewmodel.Triggers = m.Triggers

	// Load Probe
	probeType, err := probes.ProbeTypeById(m.ProbeType)
	if err != nil {
		return nil, err
	}

	probe, err := probeType.Load(m.ProbeJson)
	if err != nil {
		return nil, err
	}
	viewmodel.Probe = NewProbe(probe)

	return viewmodel, nil
}

func BlankMonitor() (*Monitor, error) {
	viewmodel := new(Monitor)

	viewmodel.Triggers = []*revere.Trigger{
		&revere.Trigger{
			TargetTemplate: targets.DefaultTargetTemplate(),
		},
	}

	viewmodel.Probe = DefaultProbe()

	return viewmodel, nil
}

func (m *Monitor) GetProbeType() probes.ProbeType {
	return m.Probe.ProbeType()
}
