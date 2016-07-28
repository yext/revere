package vm

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/util"
)

type Silence struct {
	MonitorName string
	SilenceID   db.SilenceID
	MonitorID   db.MonitorID
	Subprobes   string
	Start       time.Time
	End         time.Time
}

const (
	// TODO(fchen): fix util/time... silences sends in argument as nanoseconds, not milliseconds
	maxSilenceDuration = 14 * 24 * time.Hour
)

func (s *Silence) Id() int64 {
	return int64(s.SilenceID)
}

func NewSilence(db *db.DB, id db.SilenceID) (*Silence, error) {
	monitorSilence, err := db.LoadMonitorSilence(id)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if monitorSilence == nil {
		return nil, fmt.Errorf("Error loading silence with id: %d", id)
	}

	return newSilenceFromDB(monitorSilence), nil
}

func BlankSilence() *Silence {
	return &Silence{}
}

func newSilenceFromDB(monitorSilence *db.MonitorSilence) *Silence {
	return &Silence{
		MonitorName: monitorSilence.MonitorName,
		SilenceID:   monitorSilence.SilenceID,
		MonitorID:   monitorSilence.MonitorID,
		Subprobes:   monitorSilence.Subprobes,
		Start:       monitorSilence.Start,
		End:         monitorSilence.End,
	}
}

func AllSilences(tx *db.Tx) ([]*Silence, error) {
	monitorSilences, err := tx.LoadMonitorSilences()
	if err != nil {
		return nil, errors.Trace(err)
	}

	ss := make([]*Silence, len(monitorSilences))
	for i, monitorSilence := range monitorSilences {
		ss[i] = newSilenceFromDB(monitorSilence)
	}

	return ss, nil
}

func LoadActiveSilence(db *db.DB, id db.MonitorID, subprobe string) db.SilenceID {
	silences, err := db.LoadActiveSilencesForMonitor(id)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"monitor": id,
		}).Error("Could not load active silences. Redirecting to new silence creation.")
		return 0
	}

	for _, s := range silences {
		subprobesRegexp, err := regexp.Compile(s.Subprobes)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"silence":   s.SilenceID,
				"monitor":   s.MonitorID,
				"subprobes": s.Subprobes,
			}).Error("Could not compile silence regexp. Skipping.")
			continue
		}
		// There can by multiple silences for a subprobe, will the first matched one
		if subprobesRegexp.MatchString(subprobe) {
			return s.SilenceID
		}
	}
	return 0
}

func (s *Silence) IsCreate() bool {
	return s.Id() == 0
}

func (s *Silence) Validate(db *db.DB) (errs []string) {
	errs = append(errs, s.validate()...)
	if isCreate(s) {
		errs = append(errs, s.validateNew()...)
	} else {
		old, err := NewSilence(db, s.SilenceID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Unable to load original silence with id %d", s.SilenceID))
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
	if s.MonitorID == 0 {
		errs = append(errs, "Monitor id must be provided.")
	}

	now := time.Now()
	if (now.Sub(s.Start) > time.Minute) || now.After(s.End) {
		errs = append(errs, "Start and end must be in the future.")
	}
	return
}

func (s *Silence) validateOld(old *Silence) (errs []string) {
	if old.MonitorID != s.MonitorID {
		errs = append(errs, "Monitor name cannot be changed. Create a new silence instead.")
	}
	if old.Subprobes != s.Subprobes {
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

func (s *Silence) SetHtmlParams(values url.Values) error {
	if monitorIDStr, ok := values["monitorId"]; ok {
		if len(monitorIDStr) != 1 {
			return errors.New("Only one monitor id allowed in request")
		}

		if subprobesParams, ok := values["subprobes"]; ok {
			if len(subprobesParams) != 1 {
				return errors.New("Only one subprobe name allowed in request")
			}
			s.Subprobes = subprobesParams[0]
		}
		id, err := strconv.Atoi(monitorIDStr[0])
		if err != nil {
			return errors.Trace(err)
		}
		s.MonitorID = db.MonitorID(id)
	}
	return nil
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

func (s *Silence) Save(tx *db.Tx) error {
	monitorSilence := &db.MonitorSilence{
		MonitorName: s.MonitorName,
		Silence: &db.Silence{
			SilenceID: s.SilenceID,
			MonitorID: s.MonitorID,
			Subprobes: s.Subprobes,
			Start:     s.Start,
			End:       s.End,
		},
	}
	if isCreate(s) {
		id, err := tx.CreateMonitorSilence(monitorSilence)
		s.SilenceID = id
		return errors.Trace(err)
	} else {
		err := tx.UpdateMonitorSilence(monitorSilence)
		return errors.Trace(err)
	}
}
