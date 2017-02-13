package target

import (
	"encoding/json"

	"github.com/yext/revere/db"
)

type SlackType struct{}

type SlackTarget struct {
	SlackType
	Channel string
	// TODO(psingh): Fetch unarchived list of a channels using the APIToken provided in settings
	// This will allow the use of a dropdown on the page
}

func init() {
	addType(SlackType{})
}

func (SlackType) Id() db.TargetType {
	return 2
}

func (SlackType) Name() string {
	return "Slack"
}

func (SlackType) loadFromParams(target string) (VM, error) {
	var s SlackTarget
	err := json.Unmarshal([]byte(target), &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (SlackType) loadFromDb(encodedTarget string) (VM, error) {
	var s SlackDBModel
	err := json.Unmarshal([]byte(encodedTarget), &s)
	if err != nil {
		return nil, err
	}

	return SlackTarget{
		Channel: s.Channel,
	}, nil
}

func (SlackType) blank() VM {
	return SlackTarget{}
}

func (SlackType) Templates() map[string]string {
	return map[string]string{
		"edit": "slack-edit.html",
		"view": "slack-view.html",
	}
}

func (SlackType) Scripts() map[string][]string {
	return map[string][]string{}
}

func (st SlackTarget) Serialize() (string, error) {
	stDB := SlackDBModel{
		Channel: st.Channel,
	}

	stDBJSON, err := json.Marshal(stDB)
	return string(stDBJSON), err
}

func (SlackTarget) Type() VMType {
	return SlackType{}
}

func (et SlackTarget) Validate() (errs []string) {
	// TODO(psingh): Better validation, check channel name against slack
	if et.Channel == "" {
		errs = append(errs, "Channel name is required.")
	}
	return
}
