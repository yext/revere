package vm

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probes"
)

type Monitor struct {
	MonitorID   db.MonitorID
	Name        string
	Owner       string
	Description string
	Response    string
	ProbeType   db.ProbeType
	ProbeParams string
	Changed     time.Time
	Version     int32
	// TODO(fchen): changed and Archived need to match
	Archived *time.Time
	Probe    probes.Probe
	Triggers []*MonitorTrigger
	Labels   []*MonitorLabel
}

func (m *Monitor) Id() int64 {
	return int64(m.MonitorID)
}

func NewMonitor(tx *db.Tx, id db.MonitorID) (*Monitor, error) {
	monitor, err := tx.LoadMonitor(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if monitor == nil {
		return nil, fmt.Errorf("Monitor not found: %d", id)
	}

	m, err := newMonitorFromDB(monitor)
	if err != nil {
		return nil, errors.Trace(err)
	}

	err := m.loadComponents(tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return m, nil
}

func newMonitorFromDB(monitor *db.Monitor) (*Monitor, error) {
	var err error
	m := &Monitor{
		MonitorID:   monitor.MonitorID,
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
		return nil, errors.Trace(err)
	}

	return m, nil
}

func newMonitorsFromDB(monitors []*db.Monitor) []*Monitor {
	ms := make([]*Monitor, len(monitors))
	for i, monitor := range monitors {
		ms[i] = newMonitorFromDB(monitor)
	}
	return ms
}

func BlankMonitor() *Monitor {
	m := &Monitor{}
	m.Triggers = blankMonitorTriggers()
	m.Labels = blankMonitorLabels()
	m.Probe = probes.Default()

	return m
}

func AllMonitors(tx *db.Tx) ([]*Monitor, error) {
	monitors, err := tx.LoadMonitors()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return newMonitorsFromDB(monitors), nil
}

func AllMonitorsForLabel(tx *db.Tx, labelID db.LabelID) ([]*Monitor, error) {
	monitors, err := tx.LoadMonitorsWithLabel(labelID)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return newMonitorsFromDB(monitors), nil
}

func PopulateLabelsForMonitors(tx *db.Tx, monitors []*Monitor) error {
	mIDs := make([]db.MonitorID, len(monitors))
	for i, m := range monitors {
		mIDs[i] = m.MonitorID
	}

	mls, err := allMonitorLabels(tx, mIDs)
	if err != nil {
		return errors.Trace(err)
	}

	for _, m := range monitors {
		m.Labels = mls[m.MonitorID]
	}
	return nil
}

func (m *Monitor) loadComponents(tx *db.Tx) error {
	var err error
	m.Triggers, err = newMonitorTriggers(tx, m.MonitorID)
	if err != nil {
		return errors.Trace(err)
	}
	m.Labels, err = newMonitorLabels(tx, m.MonitorID)
	return errors.Trace(err)
}

func (m *Monitor) Validate(db *db.DB) (errs []string) {
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
		errs = append(errs, mt.validate()...)
	}

	// TODO(fchen): rethink validation between Monitors and MonitorLabels
	for _, ml := range m.Labels {
		errs = append(errs, ml.validate(db)...)
	}
	return
}

func (m *Monitor) Save(tx *db.Tx) error {
	var err error
	m.Probe, err = probes.LoadFromParams(m.ProbeType, m.ProbeParams)
	if err != nil {
		return errors.Trace(err)
	}
	monitor, err := m.toDBMonitor()
	if err != nil {
		return errors.Trace(err)
	}

	if isCreate(m) {
		m.MonitorID, err = tx.CreateMonitor(monitor)
	} else {
		err = tx.UpdateMonitor(monitor)
	}
	if err != nil {
		return errors.Trace(err)
	}

	for _, t := range m.Triggers {
		err = t.save(tx)
		if err != nil {
			return errors.Trace(err)
		}
	}

	for _, l := range m.Labels {
		err = l.save(tx)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (m *Monitor) toDBMonitor() (*db.Monitor, error) {
	probeJSON, err := m.Probe.Serialize()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &db.Monitor{
		MonitorID:   m.MonitorId,
		Name:        m.Name,
		Owner:       m.Owner,
		Description: m.Description,
		Response:    m.Response,
		ProbeType:   m.ProbeType,
		// TODO(fchen): probably wont compile as Probe.Serialize() returns a string
		// but db.Monitor.Probe expects types.JSONText
		Probe:    probeJSON,
		Changed:  m.Changed,
		Version:  m.Version,
		Archived: m.Archived,
	}, nil
}
