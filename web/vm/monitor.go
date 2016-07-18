package vm

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
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

func (m *Monitor) IsCreate() bool {
	return m.Id() == 0
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

	m, err := newMonitorFromDB(monitor, tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	err = m.loadComponents(tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return m, nil
}

func newMonitorFromDB(monitor *db.Monitor, tx *db.Tx) (*Monitor, error) {
	var err error
	m := &Monitor{
		MonitorID:   monitor.MonitorID,
		Name:        monitor.Name,
		Owner:       monitor.Owner,
		Description: monitor.Description,
		Response:    monitor.Response,
		ProbeType:   monitor.ProbeType,
		ProbeParams: "",
		Changed:     monitor.Changed,
		Version:     monitor.Version,
		Archived:    monitor.Archived,
		Probe:       nil,
		Triggers:    nil,
		Labels:      nil,
	}
	m.Probe, err = probes.LoadFromDB(monitor.ProbeType, string(monitor.Probe), tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return m, nil
}

func newMonitorsFromDB(monitors []*db.Monitor, tx *db.Tx) ([]*Monitor, error) {
	var err error
	ms := make([]*Monitor, len(monitors))
	for i, monitor := range monitors {
		ms[i], err = newMonitorFromDB(monitor, tx)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return ms, nil
}

func BlankMonitor() (*Monitor, error) {
	var err error
	m := &Monitor{}
	m.Triggers = blankMonitorTriggers()
	m.Labels = blankMonitorLabels()
	m.Probe, err = probes.Default()
	m.ProbeType = m.Probe.Id()

	return m, errors.Trace(err)
}

func AllMonitors(tx *db.Tx) ([]*Monitor, error) {
	monitors, err := tx.LoadMonitors()
	if err != nil {
		return nil, errors.Trace(err)
	}

	dbMonitors, err := newMonitorsFromDB(monitors, tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return dbMonitors, nil
}

func AllMonitorsForLabel(tx *db.Tx, labelID db.LabelID) ([]*Monitor, error) {
	DBMonitors, err := tx.LoadMonitorsWithLabel(labelID)
	if err != nil {
		return nil, errors.Trace(err)
	}

	monitors, err := newMonitorsFromDB(DBMonitors, tx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return monitors, nil
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

func (m *Monitor) Validate(DB *db.DB) (errs []string) {
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
		errs = append(errs, mt.validate(DB)...)
	}

	// TODO(fchen): rethink validation between Monitors and MonitorLabels
	for _, ml := range m.Labels {
		errs = append(errs, ml.validate(DB)...)
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
		for _, trigger := range m.Triggers {
			trigger.setMonitorID(m.MonitorID)
		}
		for _, label := range m.Labels {
			label.setMonitorID(m.MonitorID)
		}
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
		MonitorID:   m.MonitorID,
		Name:        m.Name,
		Owner:       m.Owner,
		Description: m.Description,
		Response:    m.Response,
		ProbeType:   m.ProbeType,
		// TODO(fchen): probably wont compile as Probe.Serialize() returns a string
		// but db.Monitor.Probe expects types.JSONText
		Probe:    types.JSONText(probeJSON),
		Changed:  m.Changed,
		Version:  m.Version,
		Archived: m.Archived,
	}, nil
}
