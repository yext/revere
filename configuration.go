package revere

import (
	"database/sql"
	"fmt"
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
