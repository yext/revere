package revere

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"
)

type SubprobeID int32

type Subprobe struct {
	SubprobeId  SubprobeID
	MonitorId   MonitorID
	MonitorName string
	Name        string
	Archived    *time.Time
	Status      SubprobeStatus
}

type SubprobeStatus struct {
	SubprobeId      SubprobeID
	Recorded        time.Time
	State           State
	StateStr        string
	Silenced        bool
	EnteredState    time.Time
	FmtEnteredState string
}

const allSubprobeFields = `s.subprobeid, s.monitorid, m.name as mn, s.name, s.archived,
	ss.subprobeid, ss.recorded, ss.state, ss.silenced, ss.enteredstate`

func validateSubprobe(subprobe string) (err error) {
	if _, err = regexp.Compile(subprobe); err != nil {
		return fmt.Errorf("Invalid subprobe: %s", err.Error())
	}
	return
}

func LoadSubprobesByName(db *sql.DB, monitorId MonitorID) (subprobes []*Subprobe, err error) {
	return loadSubprobes(db, fmt.Sprintf("WHERE s.monitorid = %d ORDER BY s.name", monitorId))
}

func LoadSubprobesBySeverity(db *sql.DB) (subprobes []*Subprobe, err error) {
	return loadSubprobes(db, fmt.Sprintf("WHERE ss.state != %d ORDER BY ss.state DESC, ss.enteredstate, m.name, s.name", NORMAL))
}

func LoadSubprobesBySeverityForLabel(db *sql.DB, labelId LabelID) (subprobes []*Subprobe, err error) {
	return loadSubprobes(db, fmt.Sprintf(`
		JOIN labels_monitors lm ON lm.monitorid = m.monitorid
		WHERE ss.state != %d AND lm.labelid = %d
		ORDER BY ss.state DESC, ss.enteredstate, m.name, s.name`,
		NORMAL, labelId))
}

func loadSubprobes(db *sql.DB, condition string) (subprobes []*Subprobe, err error) {
	// TODO(dp): we might need to support other orderings in the future
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM subprobes s LEFT JOIN subprobe_statuses ss ON s.subprobeid = ss.subprobeid
			JOIN monitors m ON m.monitorid = s.monitorid %s`,
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

func LoadSubprobe(db *sql.DB, subprobeId SubprobeID) (subprobe *Subprobe, err error) {
	row, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM subprobes s LEFT JOIN subprobe_statuses ss ON s.subprobeid = ss.subprobeid
			JOIN monitors m ON m.monitorid = s.monitorid
			WHERE s.subprobeid = %d
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

	var SubprobeId *SubprobeID
	var Recorded *time.Time
	var SubprobeState *int16
	var Silenced *bool
	var EnteredState *time.Time

	if err := rows.Scan(&s.SubprobeId, &s.MonitorId, &s.MonitorName, &s.Name, &s.Archived,
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
