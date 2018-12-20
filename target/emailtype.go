package target

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/setting"
	"github.com/yext/revere/state"
)

type emailType struct{}

func init() {
	registerTargetType(emailType{})
}

func (_ emailType) ID() db.TargetType {
	return 1
}

func (_ emailType) New(config types.JSONText) (Target, error) {
	return newEmail(config)
}

func (_ emailType) Alert(Db *db.DB, a *Alert, toAlert map[db.TriggerID]Target, inactive []Target) []ErrorAndTriggerIDs {
	triggerIDs := make([]db.TriggerID, 0, len(toAlert))
	for id := range toAlert {
		triggerIDs = append(triggerIDs, id)
	}

	toBuilder := newEmailListBuilder()
	replyToBuilder := newEmailListBuilder()
	for _, target := range toAlert {
		target := target.(*Email)
		toBuilder.addSlice(target.To())
		replyToBuilder.addSlice(target.ReplyTo())
	}
	for _, target := range inactive {
		target := target.(*Email)
		replyToBuilder.addSlice(target.ReplyTo())
	}
	to := toBuilder.build()
	replyTo := replyToBuilder.build()

	if len(to) == 0 {
		return nil
	}

	// TODO(eefi): Respect line length limits. Encode headers and body to
	// avoid UTF-8 causing breaks.

	emailSetting := setting.OutgoingEmailSetting{}

	dbSettings, err := Db.LoadSettingsOfType(emailSetting.Type().Id())
	if err != nil || len(dbSettings) == 0 {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "getting settings from db"),
			IDs: triggerIDs,
		}}
	}

	settingsFromDB, err := setting.LoadFromDB(emailSetting.Type().Id(), dbSettings[0].Setting)
	if err != nil {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "unparsing db settings"),
			IDs: triggerIDs,
		}}
	}

	emailSettings, found := settingsFromDB.(*setting.OutgoingEmailSetting)
	if !found {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "extracting email settings"),
			IDs: triggerIDs,
		}}
	}

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(
		"Date: %s\n", time.Now().UTC().Format(time.RFC822Z)))
	b.WriteString(fmt.Sprintf(
		"From: %s <%s>\n", emailSettings.FromName, emailSettings.FromEmail))
	b.WriteString(fmt.Sprintf(
		"Reply-To: %s\n", strings.Join(replyTo, ", ")))
	b.WriteString(fmt.Sprintf(
		"To: %s\n", strings.Join(to, ", ")))
	b.WriteString(fmt.Sprintf(
		"Subject: [%s] %s/%s\n", emailSettings.SubjectLinePrefix, a.MonitorName, a.SubprobeName))

	err = emailTmpl.Execute(&b, a)
	if err != nil {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "render email"),
			IDs: triggerIDs,
		}}
	}

	msg := []byte(strings.Replace(b.String(), "\n", "\r\n", -1))

	err = smtp.SendMail(emailSettings.SmtpServer, nil, emailSettings.FromEmail, to, msg)
	if err != nil {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "send email"),
			IDs: triggerIDs,
		}}
	}

	return nil
}

const emailText = `
{{.NewState}} is the state of {{.MonitorName}}/{{.SubprobeName}} as of {{time .Recorded}}.

{{.Host}}/monitors/{{.MonitorID}}/subprobes/{{.SubprobeID}}

{{if ne .OldState .NewState -}}
State change: {{.OldState}}->{{.NewState}}
{{- else -}}
Has been {{.NewState}} since: {{time .EnteredState}} ({{timerel .EnteredState}})
{{- end}}
{{- if not (isNormal .NewState)}}
Was last Normal at: {{time .LastNormal}} ({{timerel .LastNormal}})
{{- end}}
{{if .Description}}
Description: {{.Description}}
{{end -}}
{{if .Response}}
Suggested response: {{.Response}}
{{end -}}
{{if .Details}}
Probe reading details:

{{.Details.Text}}
{{end -}}`

var emailTmpl = template.Must(template.New("email").Funcs(template.FuncMap{
	"isNormal": func(s state.State) bool {
		return s == state.Normal
	},
	"time": func(t time.Time) string {
		return t.UTC().Format("Mon Jan 2 2006 15:04:05 MST")
	},
	"timerel": func(t time.Time) string {
		// TODO(eefi): Implement.
		return ""
	},
}).Parse(emailText))
