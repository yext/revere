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

func LoadSilences(db *sql.DB) (map[uint]Silence, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM silences s JOIN monitors m ON s.monitor_id = m.id", allSilenceFields))
	if err != nil {
		return nil, err
	}

	allSilences := make(map[uint]Silence)
	for rows.Next() {
		s, err := loadSilenceFromRow(rows)
		if err != nil {
			return nil, err
		}
		allSilences[s.Id] = *s
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
