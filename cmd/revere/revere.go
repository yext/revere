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

	"github.com/julienschmidt/httprouter"
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

	c := revere.LoadConfigs(db)
	allConfigs := &c

	router := httprouter.New()

	router.HandlerFunc("GET", "/", http.RedirectHandler("/readings/", http.StatusTemporaryRedirect).ServeHTTP)
	router.GET("/readings/", web.ReadingsIndex(db, allConfigs, &lastStates))
	router.GET("/configs/", web.ConfigsIndex(db))
	router.GET("/monitors", web.MonitorsIndex(db))
	router.GET("/monitors/:id", web.MonitorsView(db))
	router.GET("/monitors/:id/edit", web.MonitorsEdit(db))
	router.GET("/monitors/:id/subprobes", web.SubprobesIndex(db))
	router.GET("/monitors/:id/subprobes/:subprobeId", web.SubprobesView(db))
	router.GET("/configs/:id", web.ConfigsNew(db))
	router.POST("/configs/:id", web.ConfigsCreate(db))
	router.GET("/configs/:id/edit", web.ConfigsEdit(db))
	router.POST("/configs/:id/edit", web.ConfigsUpdate(db))
	router.POST("/silence", web.SilenceAlert(db))
	router.ServeFiles("/static/css/*filepath", http.Dir("web/css"))
	router.ServeFiles("/static/js/*filepath", http.Dir("web/js"))
	router.HandlerFunc("GET", "/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/favicon.ico")
	})

	go http.ListenAndServe(":"+strconv.Itoa(*port), router)

	// Initialize lastStates for the ui
	for configId, config := range *allConfigs {
		p, err := probes.NewGraphiteThreshold(config.Config)
		if err != nil {
			fmt.Printf("Config %d: Error parsing json: %s\n", configId, err.Error())
			continue
		}
		readings, err := p.Check()
		if err != nil {
			// This could be transient, so just keep going
			fmt.Printf("Config %d: Error reaching graphite: %s\n", configId, err.Error())
			continue
		}
		if lastStates[configId] == nil {
			lastStates[configId] = make(map[string]revere.State)
		}
		for subprobe, reading := range readings {
			lastStates[configId][subprobe] = reading.State
		}
	}

	ticker := time.Tick(revere.CheckFrequency * time.Minute)
	for _ = range ticker {
		c := revere.LoadConfigs(db)
		if c != nil {
			*allConfigs = c
		}
		silencedAlerts := revere.LoadSilencedAlerts(db)
		for _, config := range *allConfigs {
			// TODO(dp): validate configurations
			probe, err := probes.NewGraphiteThreshold(config.Config)
			if err != nil {
				fmt.Printf("Config %d: Error parsing json: %s\n", config.Id, err.Error())
				continue
			}

			err = runCheck(silencedAlerts, config.Id, probe, strings.Split(config.Emails, ","))
			if err != nil {
				fmt.Printf("Config %d: Error running check: %s\n", config.Id, err.Error())
			}
		}
	}
}

func runCheck(silencedAlerts map[uint]map[string]time.Time, configId uint, p *probes.GraphiteThreshold, emails []string) error {
	readings, err := p.Check()
	if err != nil {
		return err
	}
	for subprobe, reading := range readings {
		if lastStates[configId] == nil {
			lastStates[configId] = make(map[string]revere.State)
		}
		lastState := lastStates[configId][subprobe]
		lastStates[configId][subprobe] = reading.State
		shouldAlert := false
		if lastState != reading.State {
			shouldAlert = true

			// Record alert if it has changed
			_, err = db.Exec(`
			INSERT INTO readings (config_id, subprobe, state, time)
			VALUES (?, ?, ?, UTC_TIMESTAMP())
			`, configId, subprobe, reading.State)
			if err != nil {
				return err
			}
		}

		if !shouldAlert {
			shouldAlert, err = shouldSendAlert(configId, subprobe, p.AlertFrequency)
			shouldAlert = reading.State != revere.Normal && shouldAlert
			if err != nil {
				return err
			}
		}

		t := silencedAlerts[configId][subprobe]

		if shouldAlert && t.Before(time.Now()) {
			sendAlert(configId, subprobe, reading, emails)
		}
	}

	return nil
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
	lastAlert = revere.ChangeLoc(lastAlert, time.UTC)
	return !time.Now().UTC().Add(-time.Duration(alertFrequency) * time.Second).Before(lastAlert), nil
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
