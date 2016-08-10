package setting

import (
	"encoding/json"

	"github.com/yext/revere/db"
)

type OutgoingEmail struct{}

type OutgoingEmailSetting struct {
	OutgoingEmail
	FromName          string
	FromEmail         string
	SubjectLinePrefix string
	SmtpServer        string
}

// Temp struct until datasources and datasource package are combined, i.e. datasource package needs to be created
type OutgoingEmailSettingDBModel struct {
	FromName          string
	FromEmail         string
	SubjectLinePrefix string
	SmtpServer        string
}

func init() {
	addType(OutgoingEmail{})
}

func (OutgoingEmail) Id() db.SettingType {
	return 0
}

func (OutgoingEmail) Name() string {
	return "Outgoing Email Configuration"
}

func (OutgoingEmail) loadFromParams(s string) (Setting, error) {
	var oe OutgoingEmailSetting
	err := json.Unmarshal([]byte(s), &oe)
	if err != nil {
		return nil, err
	}
	return &oe, nil
}

func (OutgoingEmail) loadFromDB(s string) (Setting, error) {
	var oe OutgoingEmailSettingDBModel
	err := json.Unmarshal([]byte(s), &oe)
	if err != nil {
		return nil, err
	}

	return &OutgoingEmailSetting{
		FromName:          oe.FromName,
		FromEmail:         oe.FromEmail,
		SubjectLinePrefix: oe.SubjectLinePrefix,
		SmtpServer:        oe.SmtpServer,
	}, nil
}

func (OutgoingEmail) blank() (Setting, error) {
	return &OutgoingEmailSetting{}, nil
}

func (OutgoingEmail) Template() string {
	return "_outgoing-email.html"
}

func (OutgoingEmail) Scripts() []string {
	return []string{
		"outgoing-emails.js",
	}
}

func (oe *OutgoingEmailSetting) Serialize() (string, error) {
	oeDB := OutgoingEmailSettingDBModel{
		FromName:          oe.FromName,
		FromEmail:         oe.FromEmail,
		SubjectLinePrefix: oe.SubjectLinePrefix,
		SmtpServer:        oe.SmtpServer,
	}

	oeDBJSON, err := json.Marshal(oeDB)
	return string(oeDBJSON), err
}

func (*OutgoingEmailSetting) Type() SettingType {
	return OutgoingEmail{}
}

func (oe *OutgoingEmailSetting) Validate() []string {
	var errs []string

	if oe.FromName == "" {
		errs = append(errs, "Name is required")
	}

	if oe.FromEmail == "" {
		errs = append(errs, "Email is required")
	}

	if oe.SubjectLinePrefix == "" {
		errs = append(errs, "Subject line is required")
	}

	if oe.SmtpServer == "" {
		errs = append(errs, "SMTP server is required")
	}

	return errs
}
