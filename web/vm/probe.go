package vm

import (
	"fmt"

	"github.com/yext/revere/probes"
)

type Probe struct {
	probes.Probe
}

var (
	defaultProbeTypeId = probes.GraphiteThreshold{}.Id()
)

const (
	ProbesDir = "probes"
)

func NewProbe(p probes.Probe) *Probe {
	viewmodel := new(Probe)
	viewmodel.Probe = p

	return viewmodel
}

func DefaultProbe() *Probe {
	probe, err := BlankProbe(int(defaultProbeTypeId))
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
