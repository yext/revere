package web

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"

	"github.com/julienschmidt/httprouter"
)

type reading struct {
	Id          uint
	ConfigId    uint
	ConfigName  string
	Subprobe    string
	State       revere.State
	Time        string
	IsCurrent   bool
	SilenceTime string
}

type config struct {
	Id     uint
	Name   string
	Config string
	Emails string
	Status string
}

func ReadingsIndex(db *sql.DB, configs *map[uint]revere.Config, currentStates *map[uint]map[string]revere.State) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		silencedAlerts := revere.LoadSilencedAlerts(db)

		var readings []reading
		for configId, probeStates := range *currentStates {
			for subprobe, state := range probeStates {
				st := silencedAlerts[configId][subprobe]
				var silenceTime string
				if st.After(time.Now()) {
					silenceTime = st.Local().Format(format)
				}
				r := reading{0, configId, (*configs)[configId].Name, subprobe, state, time.Now().Format(format), true, silenceTime}
				readings = append(readings, r)
			}
		}

		rows, err := db.Query(`
		SELECT r.id, r.config_id, c.name, r.subprobe, r.state, r.time
		FROM readings r
		JOIN configurations c ON r.config_id = c.id
		ORDER BY time DESC
		`)
		if err != nil {
			fmt.Printf("Error retrieving readings: %s", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}

		for rows.Next() {
			var r reading
			var t time.Time
			if err := rows.Scan(&r.Id, &r.ConfigId, &r.ConfigName, &r.Subprobe, &r.State, &t); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err.Error())
				continue
			}
			r.Time = t.Format(format)
			st := silencedAlerts[r.ConfigId][r.Subprobe]
			if st.After(time.Now()) {
				r.SilenceTime = st.Format(format)
			}
			readings = append(readings, r)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			fmt.Printf("Got err with readings: %s\n", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}

		err = templates.ExecuteTemplate(w, "readings-index.html", map[string]interface{}{"Readings": readings})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}
	}
}

func ConfigsIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		success, err := req.Cookie("flash.success")
		flash := make(map[string][]string)
		if err == nil {
			flash["success"] = []string{success.Value}
			success.MaxAge = -1
			http.SetCookie(w, success)
		}

		c := revere.LoadConfigs(db)
		err = templates.ExecuteTemplate(w, "configs-index.html", map[string]interface{}{"Configs": c, "Flash": flash})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve configs", 500)
			return
		}
	}
}

func ConfigsNew(db *sql.DB) func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		if ps.ByName("id") != "new" {
			http.NotFound(w, req)
			return
		}

		g, err := req.Cookie("graphite")
		var graphite string
		if err == nil {
			graphite = g.Value
		}

		err = templates.ExecuteTemplate(w, "configs-new.html", map[string]interface{}{"Graphite": graphite})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to load new monitor page", 500)
			return
		}
		return
	}
}

func ParseConfigs(w http.ResponseWriter, req *http.Request, template string) (string, error) {
	if err := req.ParseForm(); err != nil {
		fmt.Println("Got err executing template:", err.Error())
		http.Error(w, "Unable to create new template", 500)
		return "", err
	}

	var errs []error
	if len(req.Form.Get("monitorName")) == 0 {
		errs = append(errs, errors.New("Monitor Name should not be empty"))
	}
	emails := req.Form.Get("emails")
	emails = strings.Replace(emails, " ", "", -1)
	if len(emails) == 0 {
		errs = append(errs, errors.New("Emails should not be empty"))
	} else {
		for _, e := range strings.Split(emails, ",") {
			if _, err := mail.ParseAddress(emails); err != nil {
				errs = append(errs, errors.New("Email: "+e+" appears invalid"))
			}
		}
	}

	config, errs2 := probes.Validate(req.Form)
	errs = append(errs, errs2...)
	if len(errs) > 0 {
		flash := make(map[string]interface{})
		flash["error"] = errs
		form := make(map[string]string)
		for k, v := range req.Form {
			form[k] = strings.Join(v, ", ")
		}
		err := templates.ExecuteTemplate(w, template, map[string]interface{}{"Flash": flash, "Form": form})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to load new monitor page", 500)
			return "", err
		}
		// Don't save bad config
		return "", errs[0]
	}
	// Replace encoded characters
	config = strings.Replace(config, "\\u003c", "<", -1)
	config = strings.Replace(config, "\\u003e", ">", -1)
	config = strings.Replace(config, "\\u0026", "&", -1)
	return config, nil
}

func ConfigsCreate(db *sql.DB) func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		if ps.ByName("id") != "new" {
			http.NotFound(w, req)
			return
		}

		config, err := ParseConfigs(w, req, "configs-new.html")
		if err != nil {
			return
		}

		_, err = db.Exec(`
		INSERT INTO configurations
		(name, config, emails)
		VALUES (?, ?, ?)
		`, req.Form.Get("monitorName"), config, req.Form.Get("emails"))
		if err != nil {
			fmt.Printf("error saving new configuration: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "flash.success",
			Value: "Successfully created monitor: " + req.Form.Get("monitorName"),
			Path:  "/configs"})
		http.SetCookie(w, &http.Cookie{
			Name:  "graphite",
			Value: req.Form.Get("graphite"),
			Path:  "/configs"})

		http.Redirect(w, req, "/configs", http.StatusFound)
	}
}

func ConfigsEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		if id == "" {
			fmt.Println("No id provided in edit configuration")
			http.Error(w, "No id provided in edit configuration", 500)
			return
		}
		var name, config, emails string
		err := db.QueryRow(
			`SELECT name, config, emails
				 FROM configurations
				 WHERE id=?
				 `, id).Scan(&name, &config, &emails)
		if err != nil {
			fmt.Println("Got err getting configuration:", err.Error())
			http.Error(w, "Unable to load monitor configuration", 500)
			return
		}

		p, err := probes.NewGraphiteThreshold(config)
		if err != nil {
			fmt.Printf("Config %s: Error parsing json %s:", id, err.Error())
			http.Error(w, "Unable to parse configuration json", 500)
			return
		}

		err = templates.ExecuteTemplate(w, "configs-edit.html", map[string]interface{}{
			"Id":     id,
			"Name":   name,
			"Emails": emails,
			"Config": p,
		})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to load edit monitor page", 500)
			return
		}
		return
	}
}

func ConfigsUpdate(db *sql.DB) func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		config, err := ParseConfigs(w, req, "configs-edit.html")
		if err != nil {
			return
		}

		_, err = db.Exec(`
		UPDATE configurations
		SET name = ?, config = ?, emails = ?
		WHERE id = ?
		`, req.Form.Get("monitorName"), config, req.Form.Get("emails"), ps.ByName("id"))
		if err != nil {
			fmt.Printf("error saving new configuration: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "flash.success",
			Value: "Successfully updated monitor: " + req.Form.Get("monitorName"),
			Path:  "/configs"})
		http.SetCookie(w, &http.Cookie{
			Name:  "graphite",
			Value: req.Form.Get("graphite"),
			Path:  "/configs"})

		http.Redirect(w, req, "/configs", http.StatusFound)
	}
}

func SilenceAlert(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		req.ParseForm()
		values := req.PostForm
		configId := values.Get("configId")
		subprobe := values.Get("subprobe")
		silenceTime, err := time.Parse(format+" -07:00", values.Get("silenceTime"))
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, "Time could not be parsed: "+values.Get("silenceTime"))
			return
		}

		if silenceTime.Before(time.Now()) {
			w.WriteHeader(400)
			io.WriteString(w, "Silence time must be in the future:"+values.Get("silenceTime"))
			return
		}

		if silenceTime.After(time.Now().Add(2 * 14 * 24 * time.Hour)) {
			w.WriteHeader(400)
			io.WriteString(w, "Silence time must be less than two weeks in the future: "+values.Get("silenceTime"))
			return
		}

		// Doing it this way assures that the database stores a UTC time
		silenceTimeUTC := revere.ChangeLoc(silenceTime.UTC(), time.Local)
		_, err = db.Exec(`
		INSERT INTO silenced_alerts
		(config_id, subprobe, silenceTime)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		silenceTime=VALUES(silenceTime)
		`, configId, subprobe, silenceTimeUTC)
		if err != nil {
			fmt.Printf("error saving alert: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	}
}
