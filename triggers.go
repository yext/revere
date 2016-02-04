package revere

import (
	"database/sql"
	"fmt"
	"html/template"
	"regexp"

	"github.com/yext/revere/targets"
	"github.com/yext/revere/util"
)

type Trigger struct {
	Id             uint                 `json:"id,omitempty"`
	Level          string               `json:"level"`
	Period         int64                `json:"period"`
	PeriodType     string               `json:"periodType"`
	Subprobe       string               `json:"subprobe"`
	Target         targets.Target       `json:"-"`
	TargetJson     string               `json:"target"`
	TargetType     targets.TargetTypeId `json:"targetType"`
	TargetTemplate template.HTML        `json:"-"`
	TriggerOnExit  bool                 `json:"triggerOnExit"`
}

const (
	allTriggerLoadFields    = "t.id, t.level, t.triggerOnExit, t.periodMs, t.targetType, t.target, mt.subprobe"
	allTriggerSaveFields    = "id, level, triggerOnExit, periodMs, targetType, target"
	allMonitorTriggerFields = "id, monitor_id, subprobe, trigger_id"
)

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
	ReverseStates = make(map[string]State)
	for k, v := range states {
		ReverseStates[v] = k
	}
}

func (t *Trigger) Validate() (errs []string) {
	if _, ok := ReverseStates[t.Level]; !ok {
		errs = append(errs, fmt.Sprintf("Invalid state for trigger: %s", t.Level))
	}

	if util.GetMs(t.Period, t.PeriodType) == 0 {
		errs = append(errs, fmt.Sprintf("Invalid period for trigger: %d %s", t.Period, t.PeriodType))
	}

	if t.Subprobe == "" {
		errs = append(errs, fmt.Sprintf("Subprobe is required"))
	}

	// Ensure subprobe is a valid regex
	if _, err := regexp.Compile(t.Subprobe); err != nil {
		errs = append(errs, fmt.Sprintf("Invalid subprobe: %s", err.Error()))
	}

	targetType, err := targets.TargetTypeById(t.TargetType)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Invalid target type for trigger: %d", t.TargetType))
	}

	target, err := targetType.Load(t.TargetJson)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Invalid target for trigger: %s", t.TargetJson))
	}
	errs = append(errs, target.Validate()...)

	return
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
	var (
		t          Trigger
		err        error
		level      State
		targetType targets.TargetType
		periodMs   int64
	)
	if err = rows.Scan(&t.Id, &level, &t.TriggerOnExit, &periodMs, &t.TargetType, &t.TargetJson, &t.Subprobe); err != nil {
		return nil, err
	}
	t.Level = States(level)
	t.Period, t.PeriodType = util.GetPeriodAndType(periodMs)

	targetType, err = targets.TargetTypeById(t.TargetType)
	if err != nil {
		return nil, err
	}

	t.Target, err = targetType.Load(t.TargetJson)
	if err != nil {
		return nil, err
	}

	t.TargetTemplate, err = t.Target.Render()
	if err != nil {
		return nil, err
	}

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

	targetType, err := targets.TargetTypeById(t.TargetType)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(nil, ReverseStates[t.Level], t.TriggerOnExit, util.GetMs(t.Period, t.PeriodType), targetType.Id(), t.TargetJson)
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
	return stmt.Close()
}

func (t *Trigger) updateTrigger(tx *sql.Tx, monitor *Monitor) (err error) {
	var stmt *sql.Stmt
	stmt, err = tx.Prepare(`UPDATE triggers t, monitor_triggers mt
		SET t.level=?, t.triggerOnExit=?, t.periodMS=?, t.targetType=?, t.target=?, mt.subprobe=?
		WHERE t.id=? AND mt.trigger_id=? AND mt.monitor_id=?`)
	if err != nil {
		return
	}

	targetType, err := targets.TargetTypeById(t.TargetType)
	if err != nil {
		return
	}

	_, err = stmt.Exec(ReverseStates[t.Level], t.TriggerOnExit, util.GetMs(t.Period, t.PeriodType), targetType.Id(), t.TargetJson, t.Subprobe, t.Id, t.Id, monitor.Id)
	if err != nil {
		return
	}
	return stmt.Close()
}
