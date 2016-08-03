package target

import (
	"encoding/json"
	"regexp"

	"github.com/yext/revere/db"
)

type EmailType struct{}

type EmailTarget struct {
	EmailType
	Addresses []*EmailAddress
}

type EmailAddress struct {
	To      string
	ReplyTo string
}

var (
	emailRegex = regexp.MustCompile(`^[\w\.\-\+\_]+@[\w\.\-]+\.[a-zA-Z]+$`)
)

func init() {
	addTargetVMType(EmailType{})
}

func (EmailType) Id() db.TargetType {
	return 1
}

func (EmailType) Name() string {
	return "Email"
}

func (EmailType) loadFromParams(target string) (TargetVM, error) {
	var e EmailTarget
	err := json.Unmarshal([]byte(target), &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (EmailType) loadFromDb(encodedTarget string) (TargetVM, error) {
	var e EmailDBModel
	err := json.Unmarshal([]byte(encodedTarget), &e)
	if err != nil {
		return nil, err
	}

	var et EmailTarget
	et.Addresses = make([]*EmailAddress, len(e.Addresses))
	for i, _ := range e.Addresses {
		et.Addresses[i] = &EmailAddress{
			To:      e.Addresses[i].To,
			ReplyTo: e.Addresses[i].ReplyTo,
		}
	}

	return et, nil
}

func (EmailType) blank() TargetVM {
	return EmailTarget{}
}

func (EmailType) Templates() map[string]string {
	return map[string]string{
		"edit": "email-edit.html",
		"view": "email-view.html",
	}
}

func (EmailType) Scripts() map[string][]string {
	return map[string][]string{
		"edit": []string{
			"email.js",
		},
	}
}

func (et EmailTarget) Serialize() (string, error) {
	etDB := EmailDBModel{}

	etDB.Addresses = make(
		[]struct {
			To      string
			ReplyTo string
		},
		len(et.Addresses),
	)

	for i, _ := range et.Addresses {
		etDB.Addresses[i].To = et.Addresses[i].To
		etDB.Addresses[i].ReplyTo = et.Addresses[i].ReplyTo
	}

	etDBJSON, err := json.Marshal(etDB)
	return string(etDBJSON), err
}

func (EmailTarget) Type() TargetTypeVM {
	return EmailType{}
}

func (et EmailTarget) Validate() (errs []string) {
	for _, e := range et.Addresses {
		if !emailRegex.MatchString(e.To) {
			errs = append(errs, "An invalid email to was provided.")
			break
		}
	}
	for _, e := range et.Addresses {
		if e.ReplyTo == "" {
			continue
		}
		if !emailRegex.MatchString(e.ReplyTo) {
			errs = append(errs, "An invalid reply to was provided.")
			break
		}
	}
	return
}
