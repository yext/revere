package vm

import (
	"fmt"

	"github.com/yext/revere/probes"
)

type Probe struct {
	probes.Probe
	templates map[string]string
	scripts   map[string][]string
}

var (
	defaultProbeType = probes.GraphiteThreshold{}.Id()
)

const (
	probesDir = "probes"
)

func NewProbe(p probes.Probe) *Probe {
	viewmodel := new(Probe)
	viewmodel.Probe = p
	probeType := p.ProbeType()
	viewmodel.templates = probeType.Templates()
	viewmodel.scripts = probeType.Scripts()

	return viewmodel
}

func DefaultProbe() *Probe {
	probe, err := BlankProbe(int(defaultProbeType))
	if err != nil {
		panic(err)
	}

	return probe
}

func BlankProbe(pt int) (*Probe, error) {
	probeType, err := probes.ProbeTypeById(probes.ProbeTypeId(pt))
	if err != nil {
		return nil, fmt.Errorf("Probe type not found: %d", pt)
	}

	probe, err := probeType.Load(`{}`)
	if err != nil {
		return nil, fmt.Errorf("Unable to load %s", probeType.Name())
	}

	return NewProbe(probe), nil
}

func (p *Probe) Templates() map[string]string {
	return p.templates
}

func (p *Probe) Scripts() map[string][]string {
	return p.scripts
}
