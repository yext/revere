package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type SilenceID int64

type Silence struct {
	SilenceId   SilenceID
	MonitorId   MonitorID
	MonitorName string
	Subprobe    string
	Start       time.Time
	End         time.Time
}

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

func (s *Silence) Update(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE silences SET monitorid=?, subprobe=?, start=?, end=? WHERE silenceid=?", s.MonitorId, s.Subprobe, s.Start, s.End, s.SilenceId)

	return err
}

func (s *Silence) Create(tx *sql.Tx) (SilenceID, error) {
	res, err := tx.Exec(
		fmt.Sprintf("INSERT INTO silences(%s) VALUES (?, ?, ?, ?, ?)", createSilenceFields),
		nil,
		s.MonitorId,
		s.Subprobe,
		s.Start,
		s.End)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return SilenceID(id), err
}
