package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/yext/revere"
)

type reading struct {
	Id         uint
	ConfigName string
	Subprobe   string
	State      revere.State
	Time       time.Time
	IsCurrent  bool
}

var (
	templates = make(map[string]*template.Template)
)

func init() {
	templates["readings-index.html"] = template.Must(template.ParseFiles("web/views/readings-index.html", "web/views/header.html", "web/views/footer.html"))
}

func ReadingsIndex(db *sql.DB, configs *map[uint]revere.Config, currentStates *map[uint]map[string]revere.State) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var readings []reading
		for configId, probeStates := range *currentStates {
			for subprobe, state := range probeStates {
				r := reading{0, (*configs)[configId].Name, subprobe, state, time.Now(), true}
				readings = append(readings, r)
			}
		}

		rRows, err := db.Query(`
		SELECT r.id, c.name, r.subprobe, r.state, r.time
		FROM readings r
		JOIN configurations c ON r.config_id = c.id
		ORDER BY time DESC
		`)
		if err != nil {
			fmt.Printf("Error retrieving readings: %s", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}

		for rRows.Next() {
			var r reading
			if err := rRows.Scan(&r.Id, &r.ConfigName, &r.Subprobe, &r.State, &r.Time); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err.Error())
				continue
			}
			readings = append(readings, r)
		}
		rRows.Close()
		if err := rRows.Err(); err != nil {
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
