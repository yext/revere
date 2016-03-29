package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
)

type Monitor struct {
	*revere.Monitor
	Probe    *Probe
	Triggers []*Trigger
	Labels   *MonitorLabels
}

func NewMonitor(db *sql.DB, id int) (*Monitor, error) {
	m, err := revere.LoadMonitor(db, uint(id))
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("Monitor not found: %d", id)
	}

	return newMonitorFromModel(db, m)
}

func newMonitorFromModel(db *sql.DB, m *revere.Monitor) (*Monitor, error) {
	viewmodel := new(Monitor)

	viewmodel.Monitor = m
	var err error
	viewmodel.Triggers, err = NewTriggersFromMonitorTriggers(m.Triggers)
	if err != nil {
		return nil, err
	}
	viewmodel.Labels, err = NewMonitorLabels(db, m.Labels)
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

func BlankMonitor(db *sql.DB) (*Monitor, error) {
	var err error
	viewmodel := new(Monitor)
	viewmodel.Monitor = new(revere.Monitor)
	viewmodel.Triggers = []*Trigger{
		BlankTrigger(),
	}
	viewmodel.Labels, err = BlankMonitorLabels(db)
	if err != nil {
		return nil, err
	}
	viewmodel.Probe = DefaultProbe()

	return viewmodel, nil
}

func AllMonitors(db *sql.DB) ([]*Monitor, error) {
	revereMonitors, err := revere.LoadMonitors(db)
	if err != nil {
		return nil, err
	}

	monitors := make([]*Monitor, len(revereMonitors))
	for i, revereMonitor := range revereMonitors {
		monitors[i], err = newMonitorFromModel(db, revereMonitor)
		if err != nil {
			return nil, err
		}
	}

	return monitors, nil
}

func (m *Monitor) GetProbeType() probes.ProbeType {
	return m.Probe.ProbeType()
}
