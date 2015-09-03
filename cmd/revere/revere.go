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
	flag.Parse()
	// TODO(dp): add more documentation
	if flag.NArg() < 3 {
		fmt.Println("Not enough arguments.\nrevere [db-hostname] [db-username] [db-password]")
		return
	}

	var (
		db_host     = flag.Arg(0)
		db_username = flag.Arg(1)
		db_password = flag.Arg(2)
	)

	dbspec := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/revere?loc=Local&parseTime=true",
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
				fmt.Printf("Error recording alert: %s\n", err.Error())
			}

			// Since state changed, we should send an alert
			sendAlert(configId, subprobe, reading, emails)
			return
		}

		shouldAlert, err := shouldSendAlert(configId, subprobe, p.AlertFrequency())
		if err != nil {
			fmt.Printf("Error checking should send alert: %s\n", err.Error())
			return
		}
		if reading.State != revere.Normal && shouldAlert {
			sendAlert(configId, subprobe, reading, emails)
		}
	}
}

func shouldSendAlert(configId uint, subprobe string, alertFrequency uint) (bool, error) {
	row := db.QueryRow(`
		SELECT time FROM alerts
		WHERE config_id = ? AND subprobe = ?
		ORDER BY time DESC
		LIMIT 1;
		`, configId, subprobe)
	var lastAlert time.Time
	err := row.Scan(&lastAlert)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return time.Now().Add(-time.Duration(alertFrequency) * time.Second).After(lastAlert), nil
}

func sendAlert(configId uint, subprobe string, reading revere.Reading, emails []string) {
	headers := make(map[string]string)
	headers["To"] = strings.Join(emails, ", ")
	if reading.State == revere.Normal {
		headers["Subject"] = "Revere reported recovery for " + subprobe
	} else {
		headers["Subject"] = "Revere reported unhealthy state for " + subprobe
	}

	b := new(bytes.Buffer)
	for k, v := range headers {
		fmt.Fprintf(b, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(b, "\r\nProbe %s reported unhealthy state with message: \n\n%s",
		subprobe, reading.Details.Text())

	err := smtp.SendMail(mailServer, nil, *sender, emails, b.Bytes())
	if err != nil {
		fmt.Printf("error sending email: %s\n", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO alerts
		(config_id, subprobe)
		VALUES (?, ?)
		`, configId, subprobe)
	if err != nil {
		fmt.Printf("error saving alert: %s\n", err.Error())
		return
	}
}
