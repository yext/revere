package target

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"
)

type Slack struct {
	Channel string
}

func newSlack(configJSON types.JSONText) (Target, error) {
	var config SlackDBModel
	err := configJSON.Unmarshal(&config)
	if err != nil {
		return nil, errors.Maskf(err, "deserialize target config")
	}

	return &Slack{Channel: config.Channel}, nil
}

func (Slack) Type() Type {
	return slackType{}
}
