package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yext/revere"
	"github.com/yext/revere/web/vm"

	"github.com/julienschmidt/httprouter"
)

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
				"Subprobes":   s,
				"MonitorName": monitorName,
				"Breadcrumbs": vm.SubprobeIndexBcs(monitorName, monitorId),
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
				"Readings":    readings,
				"Subprobe":    s,
				"MonitorName": s.MonitorName,
				"Breadcrumbs": vm.SubprobeViewBcs(s),
			})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}
