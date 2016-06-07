package vm

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/yext/revere"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probes"
)

type Monitor struct {
	MonitorId   revere.MonitorID
	Name        string
	Owner       string
	Description string
	Response    string
	ProbeType   probes.ProbeTypeId
	ProbeParams string
	Changed     time.Time
	Version     int32
	Archived    *time.Time
	Probe       probes.Probe
	Triggers    []*MonitorTrigger
	Labels      []*MonitorLabel
}

func (m *Monitor) Id() int64 {
	return int64(m.MonitorId)
}

func NewMonitor(db *sql.DB, id revere.MonitorID) (*Monitor, error) {
	monitor, err := revere.LoadMonitor(db, id)
	if err != nil {
		return nil, err
	}
	if monitor == nil {
		return nil, fmt.Errorf("Monitor not found: %d", id)
	}

	m, err := newMonitorFromModel(monitor)
	if err != nil {
		return nil, err
	}

	err := m.loadComponents(db)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func newMonitorFromModel(monitor *revere.Monitor) (*Monitor, error) {
	var err error
	m := &Monitor{
		MonitorId:   monitor.MonitorId,
		Name:        monitor.Name,
		Owner:       monitor.Owner,
		Description: monitor.Description,
		Response:    monitor.Response,
		ProbeType:   monitor.ProbeType,
		ProbeParams: nil,
		Changed:     monitor.Changed,
		Version:     monitor.Version,
		Archived:    monitor.Archived,
		Probe:       nil,
		Triggers:    nil,
		Labels:      nil,
	}
	m.Probe, err = probes.LoadFromDb(monitor.ProbeType, monitor.ProbeJson)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func newMonitorsFromModels(monitors []*revere.Monitor) []*Monitor {
	ms := make([]*Monitor, len(monitors))
	for i, monitor := range monitors {
		ms[i] = newMonitorFromModel(monitor)
	}
	return ms
}

func BlankMonitor(db *sql.DB) (*Monitor, error) {
	var err error
	m := &Monitor{}
	m.Triggers = blankMonitorTriggers()
	m.Labels = blankMonitorLabels()
	if err != nil {
		return nil, err
	}
	m.Probe = probes.Default()

	return m, nil
}

func AllMonitors(db *sql.DB) ([]*Monitor, error) {
	monitors, err := revere.LoadMonitors(db)
	if err != nil {
		return nil, err
	}

	return newMonitorsFromModels(monitors), nil
}

func AllMonitorsForLabel(db *sql.DB, labelId revere.LabelID) ([]*Monitor, error) {
	monitors, err := revere.LoadMonitorsForLabel(db, labelId)
	if err != nil {
		return nil, err
	}

	return newMonitorsFromModels(monitors), nil
}

func PopulateLabelsForMonitors(db *sql.DB, monitors []*Monitor) error {
	mIds := make([]revere.MonitorID, len(monitors))
	for i, m := range monitors {
		mIds[i] = m.MonitorId
	}

	mls, err := allMonitorLabels(db, mIds)
	if err != nil {
		return err
	}

	for _, m := range monitors {
		m.Labels = mls[m.MonitorId]
	}
	return nil
}

func (m *Monitor) loadComponents(tx *db.Tx) error {
	var err error
	m.Triggers, err = newMonitorTriggers(tx, m.MonitorId)
	if err != nil {
		return err
	}
	m.Labels, err = NewMonitorLabels(tx, m.MonitorId)
	return err
}

func (m *Monitor) Validate(db *sql.DB) (errs []string) {
	if m.Name == "" {
		errs = append(errs, fmt.Sprintf("Monitor name is required"))
	}

	var err error
	m.Probe, err = probes.LoadFromParams(m.ProbeType, m.ProbeParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load probe for monitor: %s", m.ProbeParams))
	}
	errs = append(errs, m.Probe.Validate()...)

	for _, mt := range m.Triggers {
		errs = append(errs, mt.Validate()...)
	}

	for _, ml := range m.Labels {
		errs = append(errs, ml.Validate(db)...)
	}
	return
}

func (m *Monitor) Save(tx *sql.Tx) error {
	var err error
	m.Probe, err = probes.LoadFromParams(m.ProbeType, m.ProbeParams)
	if err != nil {
		return err
	}
	probeJson, err := m.Probe.Serialize()
	if err != nil {
		return err
	}
	monitor := &revere.Monitor{
		m.MonitorId,
		m.Name,
		m.Owner,
		m.Description,
		m.Response,
		m.ProbeType,
		probeJson,
		m.Changed,
		m.Version,
		m.Archived,
	}

	if isCreate(m) {
		m.MonitorId, err = monitor.create(tx)
	} else {
		err = monitor.update(tx)
	}
	if err != nil {
		return err
	}

	for _, t := range m.Triggers {
		err = t.save(tx, m.MonitorId)
		if err != nil {
			return err
		}
	}

	for _, l := range m.Labels {
		err = l.save(tx, m.MonitorId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Monitor) toDBMonitor() (*db.Monitor, error) {
	probeJSON, err := t.Probe.Serialize()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &db.Monitor{
		MonitorID:   db.MonitorID(m.MonitorId),
		Name:        m.Name,
		Owner:       m.Owner,
		Description: m.Description,
		Response:    m.Response,
		ProbeType:   db.ProbeType(m.ProbeType),
		Probe:       probeJSON,
		Changed:     m.Changed,
		Version:     m.Version,
		Archived:    m.Archived,
	}, nil
}
