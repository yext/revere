package settings

import "encoding/json"

type OutgoingEmail struct{}

type OutgoingEmailSetting struct {
	OutgoingEmail
	FromName          string
	FromEmail         string
	SubjectLinePrefix string
	SmtpServer        string
}

func init() {
	addSettingType(OutgoingEmail{})
}

func (o OutgoingEmail) Id() SettingTypeId {
	return 0
}

func (o OutgoingEmail) Name() string {
	return "Outgoing Email Configuration"
}

func (o OutgoingEmail) Template() string {
	return "_outgoing-email.html"
}

func (o OutgoingEmail) Scripts() []string {
	return []string{
		"outgoing-emails.js",
	}
}

func (o OutgoingEmail) Load(settingJson string) (Setting, error) {
	setting := new(OutgoingEmailSetting)
	err := json.Unmarshal([]byte(settingJson), &setting)
	if err != nil {
		return nil, err
	}
	return setting, err
}

func (o OutgoingEmail) LoadDefault() Setting {
	return &OutgoingEmailSetting{}
}

func (os *OutgoingEmailSetting) Validate() []string {
	var errs []string

	if os.FromName == "" {
		errs = append(errs, "Name is required")
	}

	if os.FromEmail == "" {
		errs = append(errs, "Email is required")
	}

	if os.SubjectLinePrefix == "" {
		errs = append(errs, "Subject line is required")
	}

	if os.SmtpServer == "" {
		errs = append(errs, "SMTP server is required")
	}

	return errs
}

func (os *OutgoingEmailSetting) SettingType() SettingType {
	return OutgoingEmail{}
}
