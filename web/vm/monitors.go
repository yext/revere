package vm

import (
	"github.com/yext/revere"
	"github.com/yext/revere/probes"
)

type Monitor struct {
	*revere.Monitor
	Probe    *Probe
	Triggers []*Trigger
}

func NewMonitor(m *revere.Monitor) (*Monitor, error) {
	viewmodel := new(Monitor)

	viewmodel.Monitor = m

	var err error
	viewmodel.Triggers, err = NewTriggersFromMonitorTriggers(m.Triggers)
	if err != nil {
		return nil, err
	}

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

	viewmodel.Monitor = new(revere.Monitor)

	viewmodel.Triggers = []*Trigger{
		BlankTrigger(),
	}

	viewmodel.Probe = DefaultProbe()

	return viewmodel, nil
}

func (m *Monitor) GetProbeType() probes.ProbeType {
	return m.Probe.ProbeType()
}
