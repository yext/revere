package target

import (
	"sort"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"
)

// Email implements a target that alerts people by sending email.
type Email struct {
	to      []string
	replyTo []string
}

func newEmail(configJSON types.JSONText) (Target, error) {
	var config EmailDBModel
	err := configJSON.Unmarshal(&config)
	if err != nil {
		return nil, errors.Maskf(err, "deserialize target config")
	}

	to := newEmailListBuilder()
	replyTo := newEmailListBuilder()

	for _, a := range config.Addresses {
		to.add(a.To)
		if a.ReplyTo != "" {
			replyTo.add(a.ReplyTo)
		} else {
			replyTo.add(a.To)
		}
	}

	return &Email{to: to.build(), replyTo: replyTo.build()}, nil
}

func (e *Email) Type() Type {
	return emailType{}
}

func (e *Email) To() []string {
	return e.to[:]
}

func (e *Email) ReplyTo() []string {
	return e.replyTo[:]
}

type emailListBuilder map[string]struct{}

func newEmailListBuilder() emailListBuilder {
	return make(map[string]struct{})
}

func (b emailListBuilder) add(address string) {
	b[address] = struct{}{}
}

func (b emailListBuilder) addSlice(addresses []string) {
	for _, a := range addresses {
		b.add(a)
	}
}

func (b emailListBuilder) build() []string {
	list := make([]string, 0, len(b))
	for a := range b {
		list = append(list, a)
	}
	sort.Strings(list)
	return list
}
