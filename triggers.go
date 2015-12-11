package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type Trigger struct {
	Id            uint
	Level         string
	Period        int64
	PeriodType    string
	Subprobe      string
	Target        string
	TargetType    string
	TriggerOnExit bool
}

const allTriggerFields = "t.id, t.level, t.triggerOnExit, t.periodMs, t.targetType, t.target, mt.subprobe"

type TargetType int

const (
	Email TargetType = iota
)

var TargetTypes = map[TargetType]string{
	Email: "Email",
}

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

func ReverseStates() map[string]State {
	reverse := make(map[string]State)
	for k, v := range states {
		reverse[v] = k
	}
	return reverse
}

func States(s State) string {
	if state, ok := states[s]; ok {
		return state
	}
	return states[UNKNOWN]
}

func LoadTriggers(db *sql.DB, monitorId uint) (triggers []*Trigger, err error) {
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM triggers t JOIN monitor_triggers mt ON t.id = mt.trigger_id
				WHERE mt.monitor_id = %d
			`, allTriggerFields, monitorId))
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
