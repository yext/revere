package vm

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/util"
)

type Silence struct {
	SilenceId   revere.SilenceID
	MonitorId   revere.MonitorID
	MonitorName string
	Subprobe    string
	Start       time.Time
	End         time.Time
}

const (
	// TODO(fchen): fix util/time... silences sends in argument as nanoseconds, not milliseconds
	maxSilenceDuration = 14 * 24 * time.Hour
)

func (s *Silence) Id() int64 {
	return int64(s.SilenceId)
}

func NewSilence(db *sql.DB, id revere.SilenceID) (*Silence, error) {
	silence, err := revere.LoadSilence(db, id)
	if err != nil {
		return nil, err
	}
	if silence == nil {
		return nil, fmt.Errorf("Error loading silence with id: %d", id)
	}

	return newSilence(silence), nil
}

func BlankSilence(db *sql.DB) (*Silence, error) {
	silence := new(revere.Silence)

	return newSilence(silence), nil
}

func newSilence(s *revere.Silence) *Silence {
	return &Silence{s.SilenceId, s.MonitorId, s.MonitorName, s.Subprobe, s.Start, s.End}
}

func AllSilences(db *sql.DB) ([]*Silence, error) {
	revereSilences, err := revere.LoadSilences(db)
	if err != nil {
		return nil, err
	}

	silences := make([]*Silence, len(revereSilences))
	for i, revereSilence := range revereSilences {
		silences[i] = newSilence(revereSilence)
	}

	return silences, nil
}

func (s *Silence) Validate(db *sql.DB) (errs []string) {
	errs = append(errs, s.validate()...)
	if s.isCreate() {
		errs = append(errs, s.validateNew()...)
	} else {
		old, err := NewSilence(db, s.SilenceId)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Unable to load original silence with id %d", s.SilenceId))
		}
		errs = append(errs, s.validateOld(old)...)
	}

	return
}

func (s *Silence) validate() (errs []string) {
	if s.End.Before(s.Start) {
		errs = append(errs, "Start must be before end.")
	}

	if s.Start.Add(maxSilenceDuration).Before(s.End) {
		p, t := util.GetPeriodAndType(int64(maxSilenceDuration))
		errs = append(errs, fmt.Sprintf("End cannot be more than %d %s after start.", p, t))
	}
	return
}

func (s *Silence) validateNew() (errs []string) {
	if s.MonitorId == 0 {
		errs = append(errs, "Monitor id must be provided.")
	}

	now := time.Now()
	if (now.Sub(s.Start) > time.Minute) || now.After(s.End) {
		errs = append(errs, "Start and end must be in the future.")
	}
	return
}

func (s *Silence) validateOld(old *Silence) (errs []string) {
	if old.MonitorId != s.MonitorId {
		errs = append(errs, "Monitor name cannot be changed. Create a new silence instead.")
	}
	if old.Subprobe != s.Subprobe {
		errs = append(errs, "Subprobe cannot be changed. Create a new silence instead.")
	}

	now := time.Now()
	if old.IsPast(now) {
		return []string{"Silences from the past cannot be edited."}
	}
	if old.IsPresent(now) && !s.Start.Equal(old.Start) {
		errs = append(errs, "Start cannot be set for currently running silences.")
	}

	return
}

func (s *Silence) isCreate() bool {
	return s.SilenceId == 0
}

func (s *Silence) IsPast(moment time.Time) bool {
	return s.Start.Before(moment) && s.End.Before(moment)
}

func (s *Silence) IsPresent(moment time.Time) bool {
	return s.Start.Before(moment) && moment.Before(s.End)
}

func (s *Silence) Editable() bool {
	return time.Now().Before(s.End)
}

func (s *Silence) Save(db *sql.DB) error {
	silence := &revere.Silence{s.SilenceId, s.MonitorId, s.MonitorName, s.Subprobe, s.Start, s.End}
	if s.isCreate() {
		return revere.Transact(db, func(tx *sql.Tx) error {
			id, err := silence.Create(tx)
			s.SilenceId = id
			return err
		})
	} else {
		return revere.Transact(db, func(tx *sql.Tx) error {
			return silence.Update(tx)
		})
	}
}
