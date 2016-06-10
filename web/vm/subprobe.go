package vm

import (
	"fmt"
	"regexp"
	"time"

	"github.com/yext/revere/db"
	"github.com/yext/revere/state"
)

type SubprobeStatus struct {
	SubprobeID   db.SubprobeID
	Recorded     time.Time
	State        state.State
	Silenced     bool
	EnteredState time.Time
	LastNormal   time.Time
}

type Subprobe struct {
	SubprobeID  db.SubprobeID
	MonitorID   db.MonitorID
	MonitorName string
	Name        string
	Archived    *time.Time
	Status      SubprobeStatus
}

func (s *Subprobe) Id() int64 {
	return int64(s.SubprobeID)
}

func NewSubprobe(DB *db.DB, id db.SubprobeID) (*Subprobe, error) {
	s, err := DB.LoadSubprobe(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if s == nil {
		return nil, errors.Errorf("Subprobe not found: %d", id)
	}

	return newSubprobeFromModel(DB, s), nil
}

func newSubprobeFromDB(s *db.SubprobeWithStatusInfo) *Subprobe {
	subprobeStatus := SubprobeStatus{
		SubprobeID:   s.SubprobeID,
		Recorded:     s.Recorded,
		State:        s.State,
		Silenced:     s.Silenced,
		EnteredState: s.EnteredState,
		LastNormal:   s.LastNormal,
	}
	return &Subprobe{
		SubprobeID:  s.SubprobeID,
		MonitorID:   s.MonitorID,
		MonitorName: s.MonitorName,
		Archived:    s.Archived,
		Status:      subprobeStatus,
	}
}

func newSubprobesFromModel(ss []*db.Subprobe) []*Subprobe {
	subprobes := make([]*Subprobe, len(ss))
	for i, s := range ss {
		subprobes[i] = newSubprobeFromModel(s)
	}
	return subprobes
}

func BlankSubprobe() *Subprobe {
	return &Subprobe{}
}

func AllSubprobesFromMonitor(tx *db.Tx, id db.MonitorID) ([]*Subprobe, error) {
	ss, err := tx.LoadSubprobesByName(id)
	if err != nil {
		return nil, err
	}

	return newSubprobesFromModel(ss), nil
}

func AllAbnormalSubprobes(tx *db.Tx) ([]*Subprobe, error) {
	ss, err := tx.LoadSubprobesBySeverity()
	if err != nil {
		return nil, err
	}

	return newSubprobesFromModel(ss), nil
}

func AllAbnormalSubprobesForLabel(tx *db.Tx, id db.LabelID) ([]*Subprobe, error) {
	ss, err := tx.LoadSubprobesBySeverityForLabel(id)
	if err != nil {
		return nil, err
	}

	return newSubprobesFromModel(ss), nil
}

func AllMonitorLabelsForSubprobes(subprobes []*Subprobe) (map[db.MonitorID][]*MonitorLabel, error) {
	mIds := make([]db.MonitorID, len(subprobes))
	for i, subprobe := range subprobes {
		mIds[i] = subprobe.MonitorID
	}
	return allMonitorLabels(mIds)
}

//TODO(fchen): maybe find this a better home
func validateSubprobeRegex(subprobe string) (err error) {
	if _, err = regexp.Compile(subprobe); err != nil {
		return fmt.Errorf("Invalid subprobe: %s", err.Error())
	}
	return
}
