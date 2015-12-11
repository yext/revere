package web

// TODO(dp): rename this file once we finish migration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/yext/revere"

	"github.com/julienschmidt/httprouter"
)

const (
	format = "01/02/2006 3:04 PM"
)

var (
	tMap map[string]*template.Template
)

func init() {
	tMap = make(map[string]*template.Template)
	funcMap := template.FuncMap{"dict": dict, "isLastBc": isLastBc, "strEq": strEq}
	templateInfo, err := ioutil.ReadDir("web/views")
	for _, t := range templateInfo {
		if t.IsDir() {
			continue
		}
		tMap[t.Name()], err = template.New(t.Name()).Funcs(funcMap).ParseGlob("web/views/partials/*.html")
		if err != nil {
			panic(fmt.Sprintf("Got error initializing templates: %v", err))
		}
		tMap[t.Name()], err = tMap[t.Name()].ParseFiles("web/views/" + t.Name())
		if err != nil {
			panic(fmt.Sprintf("Got error initializing templates: %v", err))
		}
	}
}

func ActiveIssues(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		s, err := revere.LoadSubprobesBySeverity(db)
		if err != nil {
			fmt.Printf("Unable to retrieve active issues: %s", err.Error())
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues"),
				http.StatusInternalServerError)
			return
		}

		err = executeTemplate(w, "active-issues.html",
			map[string]interface{}{
				"Title":       "active issues",
				"Subprobes":   s,
				"Breadcrumbs": []breadcrumb{breadcrumb{"active issues", "/"}},
			})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
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
		err = executeTemplate(w, "monitors-index.html",
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
			fmt.Println("Got err getting triggers:", err.Error())
			http.Error(w, "Unable to retrieve monitor", http.StatusInternalServerError)
			return
		}
		err = executeTemplate(w, "monitors-view.html",
			map[string]interface{}{
				"Title":       "monitors",
				"Monitor":     m,
				"Triggers":    triggers,
				"Breadcrumbs": monitorViewBcs(m.Name, m.Id),
			})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve monitor", http.StatusInternalServerError)
			return
		}
	}
}

func MonitorsEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Monitor not found", http.StatusNotFound)
			return
		}

		// Create new monitor
		if p.ByName("id") == "new" {
			err := executeTemplate(w, "edit-monitor.html", map[string]interface{}{
				"Title": "create monitor",
			})
			if err != nil {
				fmt.Println("Unable to load new monitor page:", err.Error())
				http.Error(w, "Unable to load new monitor page", http.StatusInternalServerError)
			}
			return
		}

		// Edit existing monitor
		config := make(map[string]string)
		var configJson []byte

		i, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid monitor id: %s", err.Error()),
				http.StatusBadRequest)
			return
		}

		monitor, err := revere.LoadMonitor(db, uint(i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		triggers, err := revere.LoadTriggers(db, uint(i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		configJson, err = json.Marshal(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		err = executeTemplate(w, "edit-monitor.html", map[string]interface{}{
			"Title":    "edit monitor",
			"Monitor":  monitor,
			"Config":   string(configJson),
			"Triggers": triggers,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load edit monitor page", err.Error()),
				http.StatusInternalServerError)
		}
	}
}

func MonitorsSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// Temporarily just return success
		// TODO: Stringfy JSON in JS so that we can use a json decoder
		fmt.Println(req.FormValue("name"))
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "")
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

		s, err := revere.LoadSubprobesByName(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		var monitorName string
		var monitorId uint
		if len(s) == 0 {
			m, err := revere.LoadMonitor(db, uint(id))
			if err != nil {
				http.Error(w, "Unable to retrieve monitor", http.StatusInternalServerError)
				return
			}
			monitorName = m.Name
			monitorId = m.Id
		} else {
			monitorName = s[0].MonitorName
			monitorId = s[0].MonitorId
		}

		err = executeTemplate(w, "subprobes-index.html",
			map[string]interface{}{
				"Title":       "monitors",
				"Subprobes":   s,
				"MonitorName": monitorName,
				"Breadcrumbs": subprobeIndexBcs(monitorName, monitorId),
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

		err = executeTemplate(w, "subprobes-view.html",
			map[string]interface{}{
				"Title":       "monitors",
				"Readings":    readings,
				"Subprobe":    s,
				"MonitorName": s.MonitorName,
				"Breadcrumbs": subprobeViewBcs(s),
			})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, "Unable to retrieve subprobe", 500)
			return
		}
	}
}

func executeTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	if _, ok := data["States"]; !ok {
		data["States"] = revere.ReverseStates()
	}
	return tMap[name].ExecuteTemplate(w, name, data)
}
