package setting

import (
	"encoding/json"
	"net/url"

	"github.com/yext/revere/db"
)

type Slack struct{}

type SlackSetting struct {
	Slack
	APIToken   string
	BotName    string
	WebhookURL string
}

type SlackSettingDBModel struct {
	APIToken   string
	BotName    string
	WebhookURL string
}

func init() {
	addType(Slack{})
}

func (Slack) Id() db.SettingType {
	return 1
}

func (Slack) Name() string {
	return "Slack Configuration"
}

func (Slack) loadFromParams(s string) (Setting, error) {
	var ss SlackSetting
	err := json.Unmarshal([]byte(s), &ss)
	if err != nil {
		return nil, err
	}
	return &ss, nil
}

func (Slack) loadFromDB(s string) (Setting, error) {
	var ss SlackSettingDBModel
	err := json.Unmarshal([]byte(s), &ss)
	if err != nil {
		return nil, err
	}

	return &SlackSetting{
		APIToken:   ss.APIToken,
		BotName:    ss.BotName,
		WebhookURL: ss.WebhookURL,
	}, nil
}

func (Slack) blank() (Setting, error) {
	return &SlackSetting{}, nil
}

func (Slack) Template() string {
	return "_slack.html"
}

func (Slack) Scripts() []string {
	return []string{
		"slack.js",
	}
}

func (ss *SlackSetting) Serialize() (string, error) {
	ssDB := SlackSettingDBModel{
		APIToken:   ss.APIToken,
		BotName:    ss.BotName,
		WebhookURL: ss.WebhookURL,
	}

	ssDBJSON, err := json.Marshal(ssDB)
	return string(ssDBJSON), err
}

func (*SlackSetting) Type() SettingType {
	return Slack{}
}

func (ss *SlackSetting) Validate() []string {
	var errs []string

	// TODO(psingh): Better validation, check if valid with slack
	_, err := url.ParseRequestURI(ss.WebhookURL)
	if err != nil {
		errs = append(errs,
			"Invalid Webhook URL. Should be formatted as: "+
				"https://hooks.slack.com/services/"+
				"T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX")
	}

	if ss.APIToken == "" {
		errs = append(errs, "API Token is required")
	}
	return errs
}
