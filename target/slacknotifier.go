package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/juju/errors"
	"github.com/yext/revere/state"
)

const timeFormat = "Mon Jan 2 2006 15:04:05 MST"

var (
	stateColors = map[state.State]string{
		// Green
		state.Normal: "good",
		// Yellow
		state.Warning: "warning",
		// Red
		state.Error: "danger",
		// Black
		state.Critical: "#000000",
		// Grey
		state.Unknown: "#808080",
	}
)

type slackNotifier struct {
	alert *Alert
	name  string
	url   string
}

type payload struct {
	Username    string       `json:"username"`
	Channel     string       `json:"channel,omitempty"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Fallback  string `json:"fallback"`
	Color     string `json:"color"`
	Text      string `json:"text"`
	Timestamp int64  `json:"ts"`
}

func (s slackNotifier) sendAll(channels map[string]struct{}) error {
	var (
		failedChannelNames []string
		err                error
	)

	for channel, _ := range channels {
		err = s.send(channel)
		if err != nil {
			failedChannelNames = append(failedChannelNames, channel)
		}
	}

	if len(failedChannelNames) > 0 {
		return errors.Maskf(err, "sending slack notifications to %v", failedChannelNames)
	}
	return nil
}

func (s slackNotifier) send(channel string) error {
	message, err := s.formatMessage(channel)
	if err != nil {
		return errors.Maskf(err, "formatting slack message")
	}

	resp, err := http.Post(s.url, "application/json", message)
	if err != nil {
		return errors.Maskf(err, "sending slack notification to: %s", channel)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf(
			"not-OK HTTP status code: %d, when sending slack notification to: %s",
			resp.StatusCode,
			channel)
	}
	return nil
}

func (s slackNotifier) formatMessage(channel string) (io.Reader, error) {
	var text string
	if s.alert.OldState != s.alert.NewState {
		text = fmt.Sprintf("State change: %s->%s", s.alert.OldState, s.alert.NewState)
	} else {
		text = fmt.Sprintf("Has been %s since: %s",
			s.alert.NewState, s.alert.EnteredState.UTC().Format(timeFormat))
	}

	if s.alert.NewState != state.Normal {
		text = fmt.Sprintf("%s\nWas last Normal at: %s",
			text, s.alert.LastNormal.UTC().Format(timeFormat))
	}

	payload := payload{
		Username: s.name,
		Channel:  channel,
		Attachments: []attachment{
			{
				Title: fmt.Sprintf("%s/%s", s.alert.MonitorName, s.alert.SubprobeName),
				TitleLink: fmt.Sprintf("http://revere.khan/monitors/%d/subprobes/%d",
					s.alert.MonitorID, s.alert.SubprobeID),
				Fallback: fmt.Sprintf("%s/%s entered state: %s",
					s.alert.MonitorName, s.alert.SubprobeName, s.alert.NewState),
				Color:     stateColors[s.alert.NewState],
				Text:      text,
				Timestamp: s.alert.Recorded.Unix(),
			},
		},
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}
