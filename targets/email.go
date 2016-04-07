package targets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
)

type Email struct{}

type EmailTarget struct {
	Email          `json:"-"`
	EmailAddresses []*EmailAddress `json:"emails"`
}

type EmailAddress struct {
	EmailTo string `json:"emailTo"`
	ReplyTo string `json:"replyTo"`
}

const emailTargetTemplate = "email-edit.html"

var (
	emailTemplates = map[string]string{
		"edit": "email-edit.html",
		"view": "email-view.html",
	}
	emailScripts = map[string][]string{
		"edit": []string{
			"email.js",
		},
	}

	emailRegex = regexp.MustCompile(`^[\w\.\-\+\_]+@[\w\.\-]+\.[a-zA-Z]+$`)
)

func init() {
	addTargetType(Email{})
}

func (e Email) Id() TargetTypeId {
	return 0
}

func (e Email) Name() string {
	return "Email"
}

func (e Email) Load(target string) (Target, error) {
	var et EmailTarget
	err := json.Unmarshal([]byte(target), &et)
	if err != nil {
		return nil, err
	}
	return et, nil
}

func (et Email) Templates() map[string]string {
	return emailTemplates
}

func (et Email) Scripts() map[string][]string {
	return emailScripts
}

func (et EmailTarget) Validate() (errs []string) {
	for _, e := range et.EmailAddresses {
		if !emailRegex.MatchString(e.EmailTo) {
			errs = append(errs, "An invalid email to was provided.")
			break
		}
	}
	for _, e := range et.EmailAddresses {
		if !emailRegex.MatchString(e.ReplyTo) {
			errs = append(errs, "An invalid reply to was provided.")
			break
		}
	}
	return
}

func (et EmailTarget) TargetType() TargetType {
	return Email{}
}

func (et EmailTarget) Render() (template.HTML, error) {
	if t, ok := targetTemplates[emailTargetTemplate]; !ok {
		return template.HTML(""), fmt.Errorf("Unable to find email target template: %s", emailTargetTemplate)
	} else {
		b := bytes.Buffer{}
		t.Execute(&b, et)
		return template.HTML(b.String()), nil
	}
}
