package daemon

import (
	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
)

type trigger struct {
	*db.Trigger
	*env.Env
}

func newTrigger(dbModel *db.Trigger, env *env.Env) *trigger {
	return &trigger{Trigger: dbModel, Env: env}
}
