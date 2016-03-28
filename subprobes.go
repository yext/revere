package revere

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"
)

type Subprobe struct {
	Id          uint
	MonitorId   uint
	MonitorName string
	Name        string
	Archived    *time.Time
	Status      SubprobeStatus
}

type SubprobeStatus struct {
	SubprobeId      uint
	Recorded        time.Time
	State           State
	StateStr        string
	Silenced        bool
	EnteredState    time.Time
	FmtEnteredState string
}

const allSubprobeFields = `s.id, s.monitor_id, m.name as mn, s.name, s.archived,
	ss.subprobe_id, ss.recorded, ss.state, ss.silenced, ss.enteredState`

func validateSubprobes(subprobes string) (err error) {
	if _, err = regexp.Compile(subprobes); err != nil {
		return fmt.Errorf("Invalid subprobes: %s", err.Error())
	}
	return
}

func LoadSubprobesByName(db *sql.DB, monitorId uint) (subprobes []*Subprobe, err error) {
	return loadSubprobes(db, fmt.Sprintf("WHERE s.monitor_id = %d ORDER BY s.name", monitorId))
}

func LoadSubprobesBySeverity(db *sql.DB) (subprobes []*Subprobe, err error) {
	return loadSubprobes(db, fmt.Sprintf("WHERE ss.state != %d ORDER BY ss.state DESC, ss.enteredState, m.name, s.name", NORMAL))
}

func LoadSubprobesBySeverityWithLabel(db *sql.DB, labelId uint) (subprobes []*Subprobe, err error) {
	return loadSubprobes(db, fmt.Sprintf(`
		JOIN labels_monitors lm ON lm.monitor_id = m.id
		WHERE ss.state != %d AND lm.label_id = %d
		ORDER BY ss.state DESC, ss.enteredState, m.name, s.name`,
		NORMAL, labelId))
}

func loadSubprobes(db *sql.DB, condition string) (subprobes []*Subprobe, err error) {
	// TODO(dp): we might need to support other orderings in the future
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM subprobes s LEFT JOIN subprobe_statuses ss ON s.id = ss.subprobe_id
			JOIN monitors m ON m.id = s.monitor_id %s`,
			allSubprobeFields, condition))
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

func LoadSubprobe(db *sql.DB, subprobeId uint) (subprobe *Subprobe, err error) {
	row, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM subprobes s LEFT JOIN subprobe_statuses ss ON s.id = ss.subprobe_id
			JOIN monitors m ON m.id = s.monitor_id
			WHERE s.id = %d
			`, allSubprobeFields, subprobeId))

	if row.Next() {
		subprobe, err = loadSubprobeFromRow(row)
		if err != nil {
			return
		}
	}
	row.Close()
	if err = row.Err(); err != nil {
		return
	}
	return subprobe, nil
}

func loadSubprobeFromRow(rows *sql.Rows) (*Subprobe, error) {
	var s Subprobe

	var SubprobeId *uint
	var Recorded *time.Time
	var SubprobeState *int
	var Silenced *bool
	var EnteredState *time.Time

	if err := rows.Scan(&s.Id, &s.MonitorId, &s.MonitorName, &s.Name, &s.Archived,
		&SubprobeId, &Recorded, &SubprobeState, &Silenced, &EnteredState); err != nil {
		return nil, err
	}

	var ss SubprobeStatus
	if SubprobeId != nil {
		ss.SubprobeId = *SubprobeId
		ss.Recorded = *Recorded
		ss.State = State(*SubprobeState)
		ss.Silenced = *Silenced
		ss.EnteredState = *EnteredState
		ss.FmtEnteredState = GetFmtEnteredState(ss.EnteredState, time.Now().UTC())
	}
	ss.StateStr = States(ss.State)
	s.Status = ss

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
