package revere

import (
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Id     uint
	Name   string
	Config string
	Emails string
}

func LoadConfigs(db *sql.DB) map[uint]Config {
	cRows, err := db.Query("SELECT id, name, config, emails FROM configurations")
	if err != nil {
		fmt.Printf("Error retrieving configs: %s\n", err.Error())
		return nil
	}

	allConfigs := make(map[uint]Config)
	for cRows.Next() {
		var c Config
		if err := cRows.Scan(&c.Id, &c.Name, &c.Config, &c.Emails); err != nil {
			fmt.Printf("Error scanning rows: %s\n", err.Error())
			continue
		}
		allConfigs[c.Id] = c
	}
	cRows.Close()
	if err := cRows.Err(); err != nil {
		fmt.Printf("Got err with configs: %s\n", err.Error())
		return nil
	}

	return allConfigs
}

func LoadSilencedAlerts(db *sql.DB) map[uint]map[string]time.Time {
	silencedAlerts := make(map[uint]map[string]time.Time)
	rows, err := db.Query("SELECT sa.config_id, sa.subprobe, sa.silenceTime FROM silenced_alerts sa")
	if err != nil {
		fmt.Printf("Error retrieving silenced alerts: %s\n", err.Error())
		return nil
	}
	for rows.Next() {
		var c uint
		var sp string
		var st time.Time
		if err := rows.Scan(&c, &sp, &st); err != nil {
			fmt.Printf("Error scanning rows: %s\n", err.Error())
			continue
		}
		if silencedAlerts[c] == nil {
			silencedAlerts[c] = make(map[string]time.Time)
		}
		silencedAlerts[c][sp] = ChangeLoc(st, time.UTC)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		fmt.Printf("Error retrieving silenced alerts: %s\n", err.Error())
		return nil
	}
	return silencedAlerts
}

func ChangeLoc(t time.Time, l *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), l)
}
