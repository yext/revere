package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

type OutgoingEmailSetting struct {
	FromName          string `json:"fromName,omitempty"`
	FromEmail         string `json:"fromEmail,omitempty"`
	SubjectLinePrefix string `json:"subjectLinePrefix,omitempty"`
	SmtpServer        string `json:"smtpServer,omitempty"`
}

var outgoingEmailTemplateName = "_outgoing-email.html"

func init() {
	registerSetting(&OutgoingEmailSetting{})
}

func (m OutgoingEmailSetting) Id() int {
	return 0
}

func (m *OutgoingEmailSetting) Load() error {
	// Unimplemented
	return nil
}

func (m OutgoingEmailSetting) Save(jsonString string) error {
	// Unimplemented
	fmt.Println(jsonString)
	var testSetting OutgoingEmailSetting
	err := json.Unmarshal([]byte(jsonString), &testSetting)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v\n", testSetting)
	}
	return nil
}

func (m OutgoingEmailSetting) Render() (template.HTML, error) {
	t, ok := settingTemplates[outgoingEmailTemplateName]
	if !ok {
		return template.HTML(""), fmt.Errorf("Unable to find %s template", outgoingEmailTemplateName)
	}
	b := bytes.Buffer{}
	t.Execute(&b, m)
	return template.HTML(b.String()), nil
}
