package vm

import (
	"fmt"
	"regexp"
	"time"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/durationfmt"
	"github.com/yext/revere/state"
)

type SubprobeStatus struct {
	SubprobeID      db.SubprobeID
	Recorded        time.Time
	State           state.State
	Silenced        bool
	EnteredState    time.Time
	FmtEnteredState string
	LastNormal      time.Time
}

type Subprobe struct {
	SubprobeID  db.SubprobeID
	MonitorID   db.MonitorID
	MonitorName string
	Name        string
	EscapedName string
	Archived    *time.Time
	Status      SubprobeStatus
}

func (s *Subprobe) Id() int64 {
	return int64(s.SubprobeID)
}

func NewSubprobe(DB *db.DB, id db.SubprobeID) (*Subprobe, error) {
	s, err := DB.LoadSubprobeWithStatusInfo(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if s == nil {
		return nil, errors.Errorf("Subprobe not found: %d", id)
	}

	return newSubprobeWithStatusFromDB(s), nil
}

func newSubprobeFromDB(s *db.Subprobe) *Subprobe {
	return &Subprobe{
		SubprobeID:  s.SubprobeID,
		MonitorID:   s.MonitorID,
		MonitorName: "",
		Name:        s.Name,
		EscapedName: fmt.Sprintf("^%s$", regexp.QuoteMeta(s.Name)),
		Archived:    s.Archived,
		Status:      SubprobeStatus{},
	}
}

func newSubprobesFromDB(ss []*db.Subprobe) []*Subprobe {
	subprobes := make([]*Subprobe, len(ss))
	for i, s := range ss {
		subprobes[i] = newSubprobeFromDB(s)
	}
	return subprobes
}

func newSubprobeWithStatusFromDB(s *db.SubprobeWithStatusInfo) *Subprobe {
	subprobeStatus := SubprobeStatus{
		// TODO(fchen): maybe not hardcode in the time zone? TBD
		SubprobeID:   s.SubprobeID,
		Recorded:     s.Recorded,
		State:        s.State,
		Silenced:     s.Silenced,
		EnteredState: s.EnteredState,
		LastNormal:   s.LastNormal,
		FmtEnteredState: durationfmt.MostSigUnit().Format(
			time.Now().UTC().Sub(s.EnteredState)),
	}

	return &Subprobe{
		SubprobeID:  s.SubprobeID,
		MonitorID:   s.MonitorID,
		MonitorName: s.MonitorName,
		Name:        s.Name,
		EscapedName: fmt.Sprintf("^%s$", regexp.QuoteMeta(s.Name)),
		Archived:    s.Archived,
		Status:      subprobeStatus,
	}
}

func newSubprobesWithStatusFromDB(ss []*db.SubprobeWithStatusInfo) []*Subprobe {
	subprobesWithStatus := make([]*Subprobe, len(ss))
	for i, s := range ss {
		subprobesWithStatus[i] = newSubprobeWithStatusFromDB(s)
	}
	return subprobesWithStatus
}

func BlankSubprobe() *Subprobe {
	return &Subprobe{}
}

func AllSubprobesFromMonitor(tx *db.Tx, id db.MonitorID) ([]*Subprobe, error) {
	ss, err := tx.LoadSubprobesByName(id)
	if err != nil {
		return nil, err
	}

	return newSubprobesWithStatusFromDB(ss), nil
}

func AllAbnormalSubprobes(tx *db.Tx) ([]*Subprobe, error) {
	ss, err := tx.LoadSubprobesBySeverity()
	if err != nil {
		return nil, err
	}

	return newSubprobesWithStatusFromDB(ss), nil
}

func AllAbnormalSubprobesForLabel(tx *db.Tx, id db.LabelID) ([]*Subprobe, error) {
	ss, err := tx.LoadSubprobesBySeverityForLabel(id)
	if err != nil {
		return nil, err
	}

	return newSubprobesWithStatusFromDB(ss), nil
}

func AllMonitorLabelsForSubprobes(tx *db.Tx, subprobes []*Subprobe) (map[db.MonitorID][]*MonitorLabel, error) {
	mIds := make([]db.MonitorID, len(subprobes))
	for i, subprobe := range subprobes {
		mIds[i] = subprobe.MonitorID
	}
	return allMonitorLabels(tx, mIds)
}

//TODO(fchen): maybe find this a better home
func validateSubprobeRegex(subprobe string) (err error) {
	if _, err = regexp.Compile(subprobe); err != nil {
		return fmt.Errorf("Invalid subprobe: %s", err.Error())
	}
	return
}
