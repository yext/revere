package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
)

type Subprobe struct {
	*revere.Subprobe
}

func (s *Subprobe) Id() int64 {
	return int64(s.Subprobe.SubprobeId)
}

func (s *Subprobe) Name() string {
	return s.Subprobe.Name
}

func NewSubprobe(db *sql.DB, id revere.SubprobeID) (*Subprobe, error) {
	s, err := revere.LoadSubprobe(db, id)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, fmt.Errorf("Subprobe not found: %d", id)
	}

	return newSubprobeFromModel(db, s), nil
}

func newSubprobeFromModel(db *sql.DB, s *revere.Subprobe) *Subprobe {
	viewmodel := new(Subprobe)
	viewmodel.Subprobe = s

	return viewmodel
}

func newSubprobesFromModel(db *sql.DB, ss []*revere.Subprobe) []*Subprobe {
	subprobes := make([]*Subprobe, len(ss))
	for i, s := range ss {
		subprobes[i] = newSubprobeFromModel(db, s)
	}
	return subprobes
}

func BlankSubprobe(db *sql.DB) *Subprobe {
	viewmodel := new(Subprobe)
	viewmodel.Subprobe = new(revere.Subprobe)

	return viewmodel
}

func AllSubprobesFromMonitor(db *sql.DB, id revere.MonitorID) ([]*Subprobe, error) {
	ss, err := revere.LoadSubprobesByName(db, id)
	if err != nil {
		return nil, err
	}

	return newSubprobesFromModel(db, ss), nil
}

func AllAbnormalSubprobes(db *sql.DB) ([]*Subprobe, error) {
	ss, err := revere.LoadSubprobesBySeverity(db)
	if err != nil {
		return nil, err
	}

	return newSubprobesFromModel(db, ss), nil
}

func AllAbnormalSubprobesForLabel(db *sql.DB, id revere.LabelID) ([]*Subprobe, error) {
	ss, err := revere.LoadSubprobesBySeverityForLabel(db, id)
	if err != nil {
		return nil, err
	}

	return newSubprobesFromModel(db, ss), nil
}

func AllMonitorLabelsForSubprobes(db *sql.DB, subprobes []*Subprobe) (map[revere.MonitorID][]*MonitorLabel, error) {
	mIds := make([]revere.MonitorID, len(subprobes))
	for i, subprobe := range subprobes {
		mIds[i] = subprobe.MonitorId
	}
	return allMonitorLabels(db, mIds)
}
