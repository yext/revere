package db

import (
	"github.com/jmoiron/sqlx/types"
)

type TriggerID int32
type TargetType int16

type Trigger struct {
	TriggerID     TriggerID
	Level         State
	TriggerOnExit bool
	PeriodMilli   int32
	TargetType    TargetType
	Target        types.JSONText
}
