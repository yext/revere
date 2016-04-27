package target

import (
	"github.com/jmoiron/sqlx/types"
)

// Email implements a target that alerts people by sending email.
type Email struct {
}

func newEmail(configJSON types.JSONText) (Target, error) {
	return &Email{}, nil
}

func (e *Email) Type() Type {
	return emailType{}
}
