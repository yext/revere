/*
Revere is an alerting system for Graphite.

TODO(eefi): Detailed usage documentation.
*/
package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
)

const (
	mailServer string = "localhost:25"
)

var (
	sender *string = flag.String("emailSender", "revere@example.com", "The email from which alerts will be sent")
)

var (
	db *sql.DB

	lastStates map[uint]map[string]revere.State = make(map[uint]map[string]revere.State)
)

type config struct {
	Id     uint
	Config string
	Emails string
}

func main() {
	// TODO(dp): add more documentation
	if len(os.Args) < 4 {
		fmt.Println("Not enough arguments.\nrevere [db-hostname] [db-username] [db-password]")
		return
	}

	var (
		db_host     = os.Args[1]
		db_username = os.Args[2]
		db_password = os.Args[3]
	)

	dbspec := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/revere?loc=Local",
		db_username, db_password, db_host)

	var err error
	db, err = sql.Open("mysql", dbspec)
	if err != nil {
		fmt.Printf("Error connecting to db: %s", err.Error())
		return
	}

	fmt.Println("Loading configs from db")

	cRows, err := db.Query("SELECT * FROM configurations")

	allconfigs := make([]config, 0)
	for cRows.Next() {
		// TODO(dp): change this into debugger level logs
		var (
			id         uint
			configJson string
			email      string
		)
		err = cRows.Scan(&id, &configJson, &email)
		allconfigs = append(allconfigs, config{id, strings.Replace(configJson, "\n", "", -1), email})
		fmt.Printf("Loaded config %s\n", configJson)
	}
	cRows.Close()

	ticker := time.Tick(revere.CheckFrequency * time.Minute)
	for _ = range ticker {
		for _, config := range allconfigs {
			// TODO(dp): validate configurations
			probe, err := probes.NewGraphiteThreshold(config.Config)
			if err != nil {
				fmt.Println(err)
			}

			runCheck(config.Id, probe, strings.Split(config.Emails, ","))
		}
	}
}

func runCheck(configId uint, p *probes.GraphiteThreshold, emails []string) {
	readings, err := p.Check()
	if err != nil {
		fmt.Println(err)
		return
	}
	for subprobe, reading := range readings {
		if lastStates[configId] == nil {
			lastStates[configId] = make(map[string]revere.State)
		}
		lastState := lastStates[configId][subprobe]
		lastStates[configId][subprobe] = reading.State
		if lastState != reading.State {
			// Record alert if it has changed
			_, err = db.Exec(`
			INSERT INTO readings (config_id, subprobe, state)
			VALUES (?, ?, ?)
			`, configId, subprobe, reading.State)
			if err != nil {
				fmt.Printf("Error recording alert: %s", err.Error())
			}
		}

		if reading.State != revere.Normal {
			// Send alert
			headers := make(map[string]string)
			headers["To"] = strings.Join(emails, ", ")
			headers["Subject"] = "Revere reported unhealthy state for " + subprobe

			b := new(bytes.Buffer)
			for k, v := range headers {
				fmt.Fprintf(b, "%s: %s\r\n", k, v)
			}
			fmt.Fprintf(b, "\r\nProbe %s reported unhealthy state with message: \n\n%s",
				subprobe, reading.Details.Text())

			err = smtp.SendMail(mailServer, nil, *sender, emails, b.Bytes())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
