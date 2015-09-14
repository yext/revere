package web

import (
	"database/sql"
	"encoding/json"
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
	Id       uint
	ConfigId uint
	Config   string
	Subprobe string
	State    revere.State
	Time     time.Time
}

func ReadingsIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rRows, err := db.Query(`
		SELECT r.id, r.config_id, c.config, r.subprobe, r.state, r.time
		FROM readings r
		JOIN configurations c ON r.config_id = c.id
		`)
		if err != nil {
			fmt.Printf("Error retrieving readings: %s", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}

		var readings []reading
		for rRows.Next() {
			var r reading
			if err := rRows.Scan(&r.Id, &r.ConfigId, &r.Config, &r.Subprobe, &r.State, &r.Time); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err.Error())
				continue
			}

			// Attempt to format json
			var c interface{}
			if err := json.Unmarshal([]byte(r.Config), &c); err == nil {
				b, _ := json.MarshalIndent(c, "", "  ")
				r.Config = string(b[:])
			}
			readings = append(readings, r)
		}
		rRows.Close()
		if err := rRows.Err(); err != nil {
			fmt.Printf("Got err with readings: %s\n", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}

		t, err := template.ParseFiles("web/views/readings-index.html", "web/views/header.html", "web/views/footer.html")
		if err != nil {
			fmt.Printf("Got err parsing template: %s\n", err.Error())
			http.Error(w, "Unable to retrieve readings", 500)
			return
		}
		t.Execute(w, map[string]interface{}{"Readings": readings})
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
