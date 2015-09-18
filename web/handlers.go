package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/yext/revere"
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

const (
	format = "01/02/2006 3:04 PM"
)

var (
	templates = make(map[string]*template.Template)
)

func init() {
	templates["readings-index.html"] = template.Must(template.ParseFiles("web/views/readings-index.html", "web/views/header.html", "web/views/footer.html", "web/views/datetimepicker.html"))
	templates["configs-index.html"] = template.Must(template.ParseFiles("web/views/configs-index.html", "web/views/header.html", "web/views/footer.html"))
}

func ReadingsIndex(db *sql.DB, configs *map[uint]revere.Config, currentStates *map[uint]map[string]revere.State) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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

		err = templates["readings-index.html"].Execute(w, map[string]interface{}{"Readings": readings})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}
	}
}

func ConfigsIndex(db *sql.DB, allConfigs *map[uint]revere.Config) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		c := revere.LoadConfigs(db)
		if c != nil {
			*allConfigs = c
		}
		err := templates["configs-index.html"].Execute(w, map[string]interface{}{"Configs": *allConfigs})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve configs", 500)
			return
		}
	}
}

func SilenceAlert(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	_, filename := path.Split(r.URL.Path)
	ext := filepath.Ext(filename)
	if ext != ".css" && ext != ".js" {
		return
	}

	p := path.Join(strings.Split(r.URL.Path, "/")[2:]...)
	p = "web/" + p
	http.ServeFile(w, r, p)
}
