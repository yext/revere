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
	funcMap := template.FuncMap{"dict": dict, "lookupThreshold": lookupThreshold, "isLastBc": isLastBc}
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

		err = templates.ExecuteTemplate(w, "monitors-index.html",
			map[string]interface{}{
				"Title":       "monitors",
				"Monitors":    m,
				"Breadcrumbs": monitorIndexBcs(),
			})
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
				"Title":       "monitors",
				"Monitor":     m,
				"Triggers":    triggers,
				"Breadcrumbs": monitorViewBcs(m),
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
			map[string]interface{}{
				"Title":       "monitors",
				"Subprobes":   s,
				"Monitor":     m,
				"Breadcrumbs": subprobeIndexBcs(m),
			})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SubprobesView(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("subprobeId"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %s", p.ByName("subprobeId")),
				http.StatusNotFound)
			return
		}
		s, err := revere.LoadSubprobe(db, uint(id))
		if err != nil {
			fmt.Println("Got err getting subprobe:", err.Error())
			http.Error(w, "Unable to retrieve subprobe", http.StatusInternalServerError)
			return
		}

		if s == nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %d", id),
				http.StatusNotFound)
			return
		}

		m, err := revere.LoadMonitor(db, s.MonitorId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		readings, err := revere.LoadReadings(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve readings: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if s == nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %s", id),
				http.StatusNotFound)
			return
		}

		err = templates.ExecuteTemplate(w, "subprobes-view.html",
			map[string]interface{}{
				"Title":       "monitors",
				"Readings":    readings,
				"Subprobe":    s,
				"Monitor":     m,
				"Breadcrumbs": subprobeViewBcs(m, s),
			})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve subprobe", 500)
			return
		}
	}
}
