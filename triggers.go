package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type Trigger struct {
	Id            uint   `json:"id,omitempty"`
	Level         string `json:"level"`
	Period        int64  `json:"period"`
	PeriodType    string `json:"periodType"`
	Subprobe      string `json:"subprobe"`
	Target        string `json:"target"`
	TargetType    string `json:"targetType"`
	TriggerOnExit bool   `json:"triggerOnExit"`
}

const (
	allTriggerLoadFields    = "t.id, t.level, t.triggerOnExit, t.periodMs, t.targetType, t.target, mt.subprobe"
	allTriggerSaveFields    = "id, level, triggerOnExit, periodMs, targetType, target"
	allMonitorTriggerFields = "id, monitor_id, subprobe, trigger_id"
)

type TargetType int

const (
	Email TargetType = iota
)

var TargetTypes = map[TargetType]string{
	Email: "Email",
}

var reverseTargetTypes map[string]TargetType

type State int

const (
	NORMAL State = iota
	WARNING
	UNKNOWN
	ERROR
	CRITICAL
)

var states = map[State]string{
	NORMAL:   "NORMAL",
	UNKNOWN:  "UNKNOWN",
	WARNING:  "WARNING",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
}

var ReverseStates map[string]State

func States(s State) string {
	if state, ok := states[s]; ok {
		return state
	}
	return states[UNKNOWN]
}

func init() {
	reverseTargetTypes = make(map[string]TargetType)
	for k, v := range TargetTypes {
		reverseTargetTypes[v] = k
	}

	ReverseStates = make(map[string]State)
	for k, v := range states {
		ReverseStates[v] = k
	}
}

func LoadTriggers(db *sql.DB, monitorId uint) (triggers []*Trigger, err error) {
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM triggers t JOIN monitor_triggers mt ON t.id = mt.trigger_id
				WHERE mt.monitor_id = %d
			`, allTriggerLoadFields, monitorId))
	if err != nil {
		return nil, err
	}
	triggers = make([]*Trigger, 0)
	for rows.Next() {
		t, err := loadTriggerFromRow(rows)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, t)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return triggers, nil
}

func loadTriggerFromRow(rows *sql.Rows) (*Trigger, error) {
	var t Trigger
	var level State
	var targetType TargetType
	var periodMs int64
	var subprobe string
	if err := rows.Scan(&t.Id, &level, &t.TriggerOnExit, &periodMs, &targetType, &t.Target, &subprobe); err != nil {
		return nil, err
	}
	t.Level = States(level)
	t.TargetType = TargetTypes[targetType]
	t.Subprobe = subprobe
	t.Period, t.PeriodType = getPeriod(periodMs)
	return &t, nil
}

func (t *Trigger) saveTrigger(tx *sql.Tx, monitor *Monitor) (err error) {
	// Create/Update Trigger
	if t.Id == 0 {
		err = t.createTrigger(tx, monitor)
	} else {
		err = t.updateTrigger(tx, monitor)
	}

	return err
}

func (t *Trigger) createTrigger(tx *sql.Tx, monitor *Monitor) error {
	var stmt *sql.Stmt
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO triggers (%s) VALUES (?, ?, ?, ?, ?, ?)", allTriggerSaveFields))
	if err != nil {
		return err
	}

	res, err := stmt.Exec(nil, ReverseStates[t.Level], t.TriggerOnExit, t.getPeriodMs(), reverseTargetTypes[t.TargetType], t.Target)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}

	// Create/Update monitor_triggers
	stmt, err = tx.Prepare(fmt.Sprintf("INSERT INTO monitor_triggers(%s) VALUES (?, ?, ?, ?)", allMonitorTriggerFields))
	if err != nil {
		return err
	}

	res, err = stmt.Exec(nil, monitor.Id, t.Subprobe, id)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *Trigger) updateTrigger(tx *sql.Tx, monitor *Monitor) (err error) {
	var stmt *sql.Stmt
	stmt, err = tx.Prepare(`UPDATE triggers t, monitor_triggers mt
		SET t.level=?, t.triggerOnExit=?, t.periodMS=?, t.targetType=?, t.target=?, mt.subprobe=?
		WHERE t.id=? AND mt.trigger_id=? AND mt.monitor_id=?`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(ReverseStates[t.Level], t.TriggerOnExit, t.getPeriodMs(), reverseTargetTypes[t.TargetType], t.Target, t.Subprobe, t.Id, t.Id, monitor.Id)
	if err != nil {
		err = stmt.Close()
	}
	return
}

func getPeriod(periodMs int64) (int64, string) {
	ms := time.Duration(periodMs) * time.Millisecond
	switch {
	case ms == 0:
		return 0, ""
	case ms%(time.Hour*24) == 0:
		return int64(ms / (time.Hour * 24)), "day"
	case ms%time.Hour == 0:
		return int64(ms / time.Hour), "hour"
	case ms%time.Minute == 0:
		return int64(ms / time.Minute), "minute"
	case ms%time.Second == 0:
		return int64(ms / time.Second), "second"
	default:
		return 0, ""
	}
}

func (t *Trigger) getPeriodMs() int64 {
	switch t.PeriodType {
	case "day":
		return (t.Period * int64(time.Hour) * 24) / int64(time.Millisecond)
	case "hour":
		return (t.Period * int64(time.Hour)) / int64(time.Millisecond)
	case "minute":
		return (t.Period * int64(time.Minute)) / int64(time.Millisecond)
	case "second":
		return (t.Period * int64(time.Second)) / int64(time.Millisecond)
	default:
		return 0
	}
}
