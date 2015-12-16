package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type Silence struct {
	Id          uint
	MonitorId   uint
	MonitorName string
	Subprobes   string
	Start       time.Time
	End         time.Time
}

const allSilenceFields = "s.id, s.monitor_id, m.name, s.subprobes, s.start, s.end"

func LoadSilences(db *sql.DB) ([]*Silence, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM silences s JOIN monitors m ON s.monitor_id = m.id", allSilenceFields))
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
		`, allSilenceFields, id))
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
	if err := rows.Scan(&s.Id, &s.MonitorId, &s.MonitorName, &s.Subprobes, &s.Start, &s.End); err != nil {
		return nil, err
	}

	return &s, nil
}

// SplitSilences takes a slice of silences and returns slices of silences that have passed,
// silences that are currently in effect, and silences that will occur in the future
func SplitSilences(silences []*Silence) (past []*Silence, curr []*Silence, future []*Silence) {
	past = make([]*Silence, 0)
	curr = make([]*Silence, 0)
	future = make([]*Silence, 0)

	now := time.Now()
	for _, s := range silences {
		start := s.Start
		end := s.End
		if start.Before(now) && end.Before(now) {
			past = append(past, s)
			continue
		}
		if start.Before(now) && now.Before(end) {
			curr = append(curr, s)
			continue
		}
		// now.Before(start) && now.Before(end)
		future = append(future, s)
	}
	return
}
