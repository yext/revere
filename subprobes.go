package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type Subprobe struct {
	Id        uint
	MonitorId uint
	Name      string
	Archived  *time.Time
	Status    SubprobeStatus
}

type SubprobeStatus struct {
	SubprobeId      uint
	Recorded        time.Time
	State           Level
	StateStr        string
	Silenced        bool
	EnteredState    time.Time
	FmtEnteredState string
}

const allSubprobeFields = `s.id, s.monitor_id, s.name, s.archived,
	ss.subprobe_id, ss.recorded, ss.state, ss.silenced, ss.enteredState`

func LoadSubprobes(db *sql.DB, monitorId uint) (subprobes []*Subprobe, err error) {
	// TODO(dp): we might need to support other orderings in the future
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM subprobes s LEFT JOIN subprobe_statuses ss ON s.id = ss.subprobe_id
			WHERE s.monitor_id = %d ORDER BY s.name
			`, allSubprobeFields, monitorId))
	if err != nil {
		return nil, err
	}
	subprobes = make([]*Subprobe, 0)
	for rows.Next() {
		s, err := loadSubprobeFromRow(rows)
		if err != nil {
			return nil, err
		}
		subprobes = append(subprobes, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subprobes, nil
}
func loadSubprobeFromRow(rows *sql.Rows) (*Subprobe, error) {
	var s Subprobe

	var SubprobeId *uint
	var Recorded *time.Time
	var State *int
	var Silenced *bool
	var EnteredState *time.Time

	if err := rows.Scan(&s.Id, &s.MonitorId, &s.Name, &s.Archived,
		&SubprobeId, &Recorded, &State, &Silenced, &EnteredState); err != nil {
		return nil, err
	}

	var ss SubprobeStatus
	if SubprobeId != nil {
		ss.SubprobeId = *SubprobeId
		ss.Recorded = ChangeLoc(*Recorded, time.UTC)
		ss.State = Level(*State)
		ss.Silenced = *Silenced
		ss.EnteredState = ChangeLoc(*EnteredState, time.UTC)
		ss.FmtEnteredState = GetFmtEnteredState(ss.EnteredState, time.Now().UTC())
	}
	ss.StateStr = Levels[ss.State]
	s.Status = ss

	if s.Archived != nil {
		t := ChangeLoc(*s.Archived, time.UTC)
		s.Archived = &t
	}

	return &s, nil
}

func GetFmtEnteredState(earlier, later time.Time) string {
	d := later.Sub(earlier)

	if d <= 75*time.Minute {
		return fmt.Sprintf("%d min.", int(d.Minutes()))
	}

	if d <= 30*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}

	days := int(d.Hours()) / 24
	r := int(d.Hours()) % 24
	if r < 6 {
		return fmt.Sprintf("%d days", days)
	}
	if 6 <= r && r < 18 {
		return fmt.Sprintf("%d.5 days", days)
	}

	return fmt.Sprintf("%d days", days+1)

}
