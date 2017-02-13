package target

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/juju/errors"

	"github.com/yext/revere/db"
	"github.com/yext/revere/setting"
)

type slackType struct{}

func init() {
	registerTargetType(slackType{})
}

func (slackType) ID() db.TargetType {
	return 2
}

func (slackType) New(config types.JSONText) (Target, error) {
	return newSlack(config)
}

func (slackType) Alert(
	Db *db.DB, a *Alert, toAlert map[db.TriggerID]Target, inactive []Target) []ErrorAndTriggerIDs {
	triggerIDs := make([]db.TriggerID, 0, len(toAlert))
	for id := range toAlert {
		triggerIDs = append(triggerIDs, id)
	}

	channels := make(map[string]struct{})
	for _, target := range toAlert {
		target := target.(*Slack)
		channels[target.Channel] = struct{}{}
	}

	slackSetting := setting.SlackSetting{}
	dbSettings, err := Db.LoadSettingsOfType(slackSetting.Type().Id())
	if err != nil || len(dbSettings) == 0 {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "getting settings from db"),
			IDs: triggerIDs,
		}}
	}

	settingsFromDB, err := setting.LoadFromDB(slackSetting.Type().Id(), dbSettings[0].Setting)
	if err != nil {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "unmarshalling db settings"),
			IDs: triggerIDs,
		}}
	}

	slackSettings, found := settingsFromDB.(*setting.SlackSetting)
	if !found {
		return []ErrorAndTriggerIDs{{
			Err: errors.Maskf(err, "extracting slack settings"),
			IDs: triggerIDs,
		}}
	}

	notifier := slackNotifier{
		alert: a,
		name:  slackSettings.BotName,
		url:   slackSettings.WebhookURL,
	}
	err = notifier.sendAll(channels)
	if err != nil {
		return []ErrorAndTriggerIDs{{
			Err: errors.Trace(err),
			IDs: triggerIDs,
		}}
	}

	return nil
}
