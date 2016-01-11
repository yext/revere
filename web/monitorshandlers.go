package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yext/revere"

	"github.com/julienschmidt/httprouter"
)

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
		if p.ByName("id") == "new" {
			http.Redirect(w, req, "/monitors/new/edit", http.StatusMovedPermanently)
			return
		}

		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", p.ByName("id")), http.StatusNotFound)
			return
		}
		m, err := revere.LoadMonitor(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		if m == nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %d", id),
				http.StatusNotFound)
			return
		}

		triggers, err := revere.LoadTriggers(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func MonitorsEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", id), http.StatusNotFound)
			return
		}

		// Create new monitor
		if p.ByName("id") == "new" {
			err := executeTemplate(w, "monitors-edit.html", map[string]interface{}{
				"Title": "monitors",
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to load new monitor page: %s", err.Error()),
					http.StatusInternalServerError)
			}
			return
		}

		// Edit existing monitor
		i, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", p.ByName("id")), http.StatusNotFound)
			return
		}

		monitor, err := revere.LoadMonitor(db, uint(i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		if monitor == nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %d", i), http.StatusNotFound)
			return
		}

		triggers, err := revere.LoadTriggers(db, uint(i))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		err = executeTemplate(w, "monitors-edit.html", map[string]interface{}{
			"Title":    "monitors",
			"Monitor":  monitor,
			"Triggers": triggers,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load edit monitor page: %s", err.Error()),
				http.StatusInternalServerError)
		}
	}
}

func MonitorsSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var m *revere.Monitor
		err := json.NewDecoder(req.Body).Decode(&m)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		errs := m.Validate()
		if errs != nil {
			errors, err := json.Marshal(map[string][]string{"errors": errs})
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(errors)
			return
		}
		err = m.SaveMonitor(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		redirect, err := json.Marshal(map[string]string{"redirect": fmt.Sprintf("/monitors/%d", m.Id)})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(redirect)
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
				http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
					http.StatusInternalServerError)
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
		mId, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", p.ByName("id")),
				http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(p.ByName("subprobeId"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %s", p.ByName("subprobeId")),
				http.StatusNotFound)
			return
		}

		s, err := revere.LoadSubprobe(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		if s == nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %d", id),
				http.StatusNotFound)
			return
		}

		if s.MonitorId != uint(mId) {
			http.Error(w, fmt.Sprintf("Subprobe %d does not exist for monitor: %d", id, mId),
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
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}
