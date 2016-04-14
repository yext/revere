package revere

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yext/revere/util"
)

type SilenceID int64

type Silence struct {
	SilenceId   SilenceID `json:"id",omitempty`
	MonitorId   MonitorID `json:"monitor"`
	MonitorName string    `json:"-"`
	Subprobe    string    `json:"subprobe"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

const silenceEndLimit = 14 * 24 * time.Hour

const loadSilenceFields = "s.silenceid, s.monitorid, s.subprobe, s.start, s.end, m.name"
const createSilenceFields = "silenceid, monitorid, subprobe, start, end"

func LoadSilences(db *sql.DB) ([]*Silence, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM silences s JOIN monitors m ON s.monitorid = m.monitorid", loadSilenceFields))
	if err != nil {
		return nil, err
	}

	allSilences := make([]*Silence, 0)
	for rows.Next() {
		s, err := loadSilenceFromRow(rows)
		if err != nil {
			return nil, err
		}
		allSilences = append(allSilences, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allSilences, nil
}

func LoadSilence(db *sql.DB, id SilenceID) (s *Silence, err error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM silences s JOIN monitors m ON s.monitorid = m.monitorid WHERE s.silenceid = %d
		`, loadSilenceFields, id))
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		s, err = loadSilenceFromRow(rows)
		if err != nil {
			return nil, err
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

func loadSilenceFromRow(rows *sql.Rows) (*Silence, error) {
	var s Silence
	if err := rows.Scan(&s.SilenceId, &s.MonitorId, &s.Subprobe, &s.Start, &s.End, &s.MonitorName); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Silence) Save(db *sql.DB) error {
	if s.IsCreate() {
		id, err := s.create(db)
		if err != nil {
			return err
		}
		s.SilenceId = id
		return err
	}

	err := s.update(db)
	return err
}

func (s *Silence) update(db *sql.DB) error {
	_, err := db.Exec("UPDATE silences SET monitorid=?, subprobe=?, start=?, end=? WHERE silenceid=?",
		s.MonitorId, s.Subprobe, s.Start, s.End, s.SilenceId)
	return err
}

func (s *Silence) create(db *sql.DB) (SilenceID, error) {
	res, err := db.Exec(fmt.Sprintf(`
		INSERT INTO silences(%s) VALUES (?, ?, ?, ?, ?)
		`, createSilenceFields),
		nil, s.MonitorId, s.Subprobe, s.Start, s.End)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return SilenceID(id), err
}

func (s *Silence) ValidateAgainstOld(oldS *Silence) (errs []string) {
	if oldS == nil {
		oldS = new(Silence)
	}

	now := time.Now()
	if !oldS.IsCreate() && oldS.IsPast(now) {
		return []string{"Silences from the past cannot be edited."}
	}

	if oldS.IsCreate() {
		errs = append(errs, s.validateNewParams(now)...)
	} else {
		errs = append(errs, s.validateUpdateParams(oldS)...)
	}

	if s.End.Before(s.Start) {
		errs = append(errs, "Start must be before end.")
	}

	if s.Start.Add(silenceEndLimit).Before(s.End) {
		p, t := util.GetPeriodAndType(int64(silenceEndLimit))
		errs = append(errs, fmt.Sprintf("End cannot be more than %d %s after start.", p, t))
	}

	if oldS.IsPresent(now) && !s.Start.Equal(oldS.Start) {
		errs = append(errs, "Start cannot be set for currently running silences.")
	}

	return errs
}

func (newS *Silence) validateNewParams(now time.Time) (errs []string) {
	if newS.MonitorId == 0 {
		errs = append(errs, "Monitor id must be provided.")
	}

	if now.After(newS.Start) && now.After(newS.End) {
		errs = append(errs, "Start and end must be in the future.")
	}

	return
}

func (updateS *Silence) validateUpdateParams(oldS *Silence) (errs []string) {
	if oldS.MonitorId != updateS.MonitorId {
		errs = append(errs, "Monitor name cannot be changed. Create a new silence instead.")
	}
	if oldS.Subprobe != updateS.Subprobe {
		errs = append(errs, "Subprobe cannot be changed. Create a new silence instead.")
	}

	return
}

func (s Silence) IsCreate() bool {
	return s.SilenceId == 0
}

func (s Silence) IsPast(now time.Time) bool {
	return s.Start.Before(now) && s.End.Before(now)
}

func (s Silence) IsPresent(now time.Time) bool {
	return s.Start.Before(now) && now.Before(s.End)
}

func (s Silence) Editable() bool {
	return time.Now().Before(s.End)
}
