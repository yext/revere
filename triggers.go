package revere

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere/targets"
	"github.com/yext/revere/util"
)

type Trigger struct {
	Id            uint                 `json:"id,omitempty"`
	Level         string               `json:"level"`
	Period        int64                `json:"period"`
	PeriodType    string               `json:"periodType"`
	TargetJson    string               `json:"target"`
	TargetType    targets.TargetTypeId `json:"targetType"`
	TriggerOnExit bool                 `json:"triggerOnExit"`
}

const (
	allTriggerLoadFields = "t.id, t.level, t.triggerOnExit, t.periodMs, t.targetType, t.target"
	allTriggerSaveFields = "id, level, triggerOnExit, periodMs, targetType, target"
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

func LoadMonitorTriggers(db *sql.DB, monitorId uint) (triggers []*MonitorTrigger, err error) {
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s, mt.subprobe FROM triggers t JOIN monitor_triggers mt ON t.id = mt.trigger_id
				WHERE mt.monitor_id = %d
			`, allTriggerLoadFields, monitorId))
	if err != nil {
		return nil, err
	}
	triggers = make([]*MonitorTrigger, 0)
	for rows.Next() {
		t, err := loadMonitorTriggerFromRow(rows)
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

func loadMonitorTriggerFromRow(rows *sql.Rows) (*MonitorTrigger, error) {
	var (
		t          MonitorTrigger
		err        error
		level      State
		targetType targets.TargetType
		periodMs   int64
	)
	if err = rows.Scan(&t.Id, &level, &t.TriggerOnExit, &periodMs, &t.TargetType, &t.TargetJson, &t.Subprobes); err != nil {
		return nil, err
	}
	//TODO(psingh): Move into view monitor
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

	return &t, nil
}

func LoadLabelTriggers(db *sql.DB, labelId uint) (triggers []*LabelTrigger, err error) {
	rows, err := db.Query(
		fmt.Sprintf(`
			SELECT %s FROM triggers t JOIN label_triggers lt ON t.id = lt.trigger_id
				WHERE lt.label_id = %d
			`, allTriggerLoadFields, labelId))
	if err != nil {
		return nil, err
	}
	triggers = make([]*LabelTrigger, 0)
	for rows.Next() {
		t, err := loadLabelTriggerFromRow(rows)
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

func loadLabelTriggerFromRow(rows *sql.Rows) (*LabelTrigger, error) {
	var (
		t          LabelTrigger
		err        error
		level      State
		targetType targets.TargetType
		periodMs   int64
	)
	if err = rows.Scan(&t.Id, &level, &t.TriggerOnExit, &periodMs, &t.TargetType, &t.TargetJson); err != nil {
		return nil, err
	}
	//TODO(psingh): Move into view monitor
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

	return &t, nil
}

func (t *Trigger) Validate() (errs []string) {
	if _, ok := ReverseStates[t.Level]; !ok {
		errs = append(errs, fmt.Sprintf("Invalid state for trigger: %s", t.Level))
	}

	if util.GetMs(t.Period, t.PeriodType) == 0 {
		errs = append(errs, fmt.Sprintf("Invalid period for trigger: %d %s", t.Period, t.PeriodType))
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

func (t *Trigger) save(tx *sql.Tx) (newId uint, err error) {
	// Create/Update Trigger
	if t.Id == 0 {
		newId, err = t.create(tx)
	} else {
		err = t.update(tx)
	}

	return newId, err
}

func (t *Trigger) create(tx *sql.Tx) (uint, error) {
	var stmt *sql.Stmt
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO triggers (%s) VALUES (?, ?, ?, ?, ?, ?)", allTriggerSaveFields))
	if err != nil {
		return 0, err
	}

	targetType, err := targets.TargetTypeById(t.TargetType)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(nil, ReverseStates[t.Level], t.TriggerOnExit, util.GetMs(t.Period, t.PeriodType), targetType.Id(), t.TargetJson)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), stmt.Close()
}

func (t *Trigger) update(tx *sql.Tx) (err error) {
	var stmt *sql.Stmt
	stmt, err = tx.Prepare(`
		UPDATE triggers
		SET level=?, triggerOnExit=?, periodMS=?, targetType=?, target=?
		WHERE id=?
	`)
	if err != nil {
		return
	}

	targetType, err := targets.TargetTypeById(t.TargetType)
	if err != nil {
		return
	}

	_, err = stmt.Exec(ReverseStates[t.Level], t.TriggerOnExit, util.GetMs(t.Period, t.PeriodType), targetType.Id(), t.TargetJson, t.Id)
	if err != nil {
		return
	}
	return stmt.Close()
}

func (t *Trigger) delete(tx *sql.Tx) (err error) {
	var stmt *sql.Stmt
	stmt, err = tx.Prepare(`
		DELETE FROM triggers
		WHERE id=?
	`)
	if err != nil {
		return
	}

	_, err = stmt.Exec(t.Id)
	if err != nil {
		return
	}
	return stmt.Close()
}
