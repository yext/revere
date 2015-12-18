package revere

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yext/revere/util"
)

type Silence struct {
	Id          uint      `json:"-"`
	MonitorId   uint      `json:"monitor"`
	MonitorName string    `json:"-"`
	Subprobes   string    `json:"subprobes"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

const silenceEndLimit = 14 * 24 * time.Hour

const loadSilenceFields = "s.id, s.monitor_id, s.subprobes, s.start, s.end, m.name"
const createSilenceFields = "id, monitor_id, subprobes, start, end"

func LoadSilences(db *sql.DB) ([]*Silence, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM silences s JOIN monitors m ON s.monitor_id = m.id", loadSilenceFields))
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

func LoadSilence(db *sql.DB, id uint) (s *Silence, err error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM silences s JOIN monitors m ON s.monitor_id = m.id WHERE s.id = %d
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
	if err := rows.Scan(&s.Id, &s.MonitorId, &s.Subprobes, &s.Start, &s.End, &s.MonitorName); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Silence) Save(db *sql.DB) error {
	if s.IsNew() {
		id, err := s.create(db)
		if err != nil {
			return err
		}
		s.Id = id
		return err
	}

	err := s.update(db)
	return err
}

func (s *Silence) update(db *sql.DB) error {
	_, err := db.Exec("UPDATE silences SET monitor_id=?, subprobes=?, start=?, end=? WHERE id=?",
		s.MonitorId, s.Subprobes, s.Start, s.End, s.Id)
	return err
}

func (s *Silence) create(db *sql.DB) (uint, error) {
	res, err := db.Exec(fmt.Sprintf(`
		INSERT INTO silences(%s) VALUES (?, ?, ?, ?, ?)
		`, createSilenceFields),
		nil, s.MonitorId, s.Subprobes, s.Start, s.End)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint(id), err
}

func (s *Silence) Validate(db *sql.DB) (errs []string, err error) {
	oldS, err := LoadSilence(db, s.Id)
	if oldS == nil {
		oldS = new(Silence)
	}

	now := time.Now()
	if !oldS.IsNew() && oldS.IsPast(now) {
		return []string{"Silences from the past cannot be edited."}, nil
	}

	if oldS.IsNew() {
		errs = append(errs, oldS.validateNewParams(s, now)...)
	} else {
		errs = append(errs, oldS.validateUpdateParams(s)...)
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

	return errs, nil
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
	if oldS.Subprobes != updateS.Subprobes {
		errs = append(errs, "Subprobes cannot be changed. Create a new silence instead.")
	}

	return
}

// SplitSilences takes a slice of silences and returns slices of silences that have passed,
// silences that are currently in effect, and silences that will occur in the future
func SplitSilences(silences []*Silence) (past []*Silence, curr []*Silence, future []*Silence) {
	past = make([]*Silence, 0)
	curr = make([]*Silence, 0)
	future = make([]*Silence, 0)

	now := time.Now()
	for _, s := range silences {
		if s.IsPast(now) {
			past = append(past, s)
			continue
		}
		if s.IsPresent(now) {
			curr = append(curr, s)
			continue
		}
		// now.Before(start) && now.Before(end)
		future = append(future, s)
	}
	return
}

func (s Silence) IsNew() bool {
	return s.Id == 0
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
