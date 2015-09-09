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
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/web"
)

const (
	mailServer string = "localhost:25"
)

var (
	sender *string = flag.String("emailSender", "revere@example.com", "The email from which alerts will be sent")

	port *int = flag.Int("port", 8123, "The port on which revere will listen")
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
		dbHost     = flag.Arg(0)
		dbUsername = flag.Arg(1)
		dbPassword = flag.Arg(2)

		dbspec = fmt.Sprintf(
			"%s:%s@tcp(%s:3306)/revere?loc=Local&parseTime=true",
			dbUsername, dbPassword, dbHost)
	)

	var err error
	db, err = sql.Open("mysql", dbspec)
	if err != nil {
		fmt.Printf("Error connecting to db: %s", err.Error())
		return
	}

	http.HandleFunc("/", web.ReadingsIndex(db))
	http.HandleFunc("/static/", web.StaticHandler)

	go http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	allConfigs := loadConfigs()

	ticker := time.Tick(revere.CheckFrequency * time.Minute)
	for _ = range ticker {
		for _, config := range allConfigs {
			// TODO(dp): validate configurations
			probe, err := probes.NewGraphiteThreshold(config.Config)
			if err != nil {
				fmt.Println("Error parsing json:", err.Error())
				continue
			}

			runCheck(config.Id, probe, strings.Split(config.Emails, ","))
		}
	}
}

func loadConfigs() []config {
	fmt.Println("Loading configs from db")

	cRows, err := db.Query("SELECT * FROM configurations")
	if err != nil {
		fmt.Printf("Error retrieving configs: %s", err.Error())
		return nil
	}

	var allConfigs []config
	for cRows.Next() {
		// TODO(dp): change this into debugger level logs
		var c config
		if err := cRows.Scan(&c.Id, &c.Config, &c.Emails); err != nil {
			fmt.Printf("Error scanning rows: %s\n", err.Error())
			continue
		}
		allConfigs = append(allConfigs, c)
		fmt.Printf("Loaded config %s\n", c.Config)
	}
	cRows.Close()
	if err := cRows.Err(); err != nil {
		fmt.Printf("Got err with configs: %s\n", err.Error())
		return nil
	}

	return allConfigs
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
			INSERT INTO readings (config_id, subprobe, state, time)
			VALUES (?, ?, ?, UTC_TIMESTAMP())
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
	var stateString string
	if reading.State == revere.Normal {
		stateString = "recovery"
	} else {
		stateString = "unhealthy"
	}

	if reading.State == revere.Normal {
		headers["Subject"] = "Revere reported " + stateString + " for " + subprobe
	} else {
		headers["Subject"] = "Revere reported " + stateString + " for " + subprobe
	}

	b := new(bytes.Buffer)
	for k, v := range headers {
		fmt.Fprintf(b, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(b, "\r\nProbe %s reported %s state with message: \n\n%s",
		subprobe, stateString, reading.Details.Text())

	err := smtp.SendMail(mailServer, nil, *sender, emails, b.Bytes())
	if err != nil {
		fmt.Printf("error sending email: %s\n", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO alerts
		(config_id, subprobe, time)
		VALUES (?, ?, UTC_TIMESTAMP())
		`, configId, subprobe)
	if err != nil {
		fmt.Printf("error saving alert: %s\n", err.Error())
		return
	}
}
