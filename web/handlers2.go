package web

// TODO(dp): rename this file once we finish migration

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/yext/revere"

	"github.com/julienschmidt/httprouter"
)

const (
	format = "01/02/2006 3:04 PM"
)

var (
	templates *template.Template
)

func init() {
	funcMap := template.FuncMap{"dict": dict, "lookupThreshold": lookupThreshold}
	var err error
	templates, err = template.New("").Funcs(funcMap).ParseGlob("web/views/*.html")
	if err != nil {
		panic(fmt.Sprintf("Got error initializing templates: %v", err))
	}
}

func MonitorsIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		m, err := revere.LoadMonitors(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitors: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = templates.ExecuteTemplate(w, "monitors-index.html", map[string]interface{}{"Monitors": m})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitors: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func MonitorsView(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.NotFound(w, req)
			return
		}
		m, err := revere.LoadMonitor(db, uint(id))
		if err != nil {
			fmt.Println("Got err getting monitor:", err.Error())
			http.Error(w, "Unable to retrieve monitor", http.StatusInternalServerError)
			return
		}
		triggers, err := revere.LoadTriggers(db, uint(id))
		if err != nil {
			fmt.Println("Got err getting monitor:", err.Error())
			http.Error(w, "Unable to retrieve monitor", http.StatusInternalServerError)
			return
		}
		err = templates.ExecuteTemplate(w, "monitors-view.html",
			map[string]interface{}{
				"Monitor":  m,
				"Triggers": triggers,
			})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve monitor", 500)
			return
		}
	}
}

func MonitorsEdit(_ *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		http.Error(w, fmt.Sprintf("Not yet implemented for id %s", p.ByName("id")), http.StatusNotImplemented)
	}
}

func SubprobesIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", p.ByName("id")),
				http.StatusNotFound)
			return
		}

		s, err := revere.LoadSubprobes(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		m, err := revere.LoadMonitor(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		if m == nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %d", id),
				http.StatusNotFound)
			return
		}

		err = templates.ExecuteTemplate(w, "subprobes-index.html",
			map[string]interface{}{"Subprobes": s, "MonitorName": m.Name})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SubprobesView(_ *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		http.Error(w, fmt.Sprintf("Not yet implemented for subprobeId %s", p.ByName("subprobeId")), http.StatusNotImplemented)
	}
}
