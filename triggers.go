package revere

import (
	"database/sql"
	"fmt"
)

type Trigger struct {
	Id            uint
	Subprobe      string
	Level         string
	TriggerOnExit bool
	PeriodMs      int
	TargetType    string
	Target        string
}

const allTriggerFields = "t.id, t.level, t.triggerOnExit, t.periodMs, t.targetType, t.target, mt.subprobe"

type TargetType int

const (
	Email TargetType = iota
)

var TargetTypes = map[TargetType]string{
	Email: "Email",
}

type Level int

const (
	UNKNOWN Level = iota
	NORMAL
	WARNING
	ERROR
	CRITICAL
)

var Levels = map[Level]string{
	NORMAL:   "NORMAL",
	WARNING:  "WARNING",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
	UNKNOWN:  "UNKNOWN",
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
	var level Level
	var targetType TargetType
	var subprobe string
	if err := rows.Scan(&t.Id, &level, &t.TriggerOnExit, &t.PeriodMs, &targetType, &t.Target, &subprobe); err != nil {
		return nil, err
	}
	t.Level = Levels[level]
	t.TargetType = TargetTypes[targetType]
	t.Subprobe = subprobe
	return &t, nil
}
