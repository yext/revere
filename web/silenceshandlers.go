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

func SilencesIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		s, err := revere.LoadSilences(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silences: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		past, curr, future := revere.SplitSilences(s)

		err = executeTemplate(w, "silences-index.html",
			silenceDataWith(map[string]interface{}{
				"Past":        past,
				"Curr":        curr,
				"Future":      future,
				"Breadcrumbs": silencesIndexBcs(),
			}))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silences: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SilencesView(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.NotFound(w, req)
			return
		}
		s, err := revere.LoadSilence(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if s == nil {
			http.Error(w, fmt.Sprintf("Silence not found: %d", id),
				http.StatusNotFound)
			return
		}

		err = executeTemplate(w, "silences-view.html",
			silenceDataWith(map[string]interface{}{
				"Silence":     s,
				"Breadcrumbs": silencesViewBcs(s.Id, s.MonitorName),
			}))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func SilencesEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		idStr := p.ByName("id")
		if idStr == "" {
			http.Error(w, fmt.Sprintf("Silence not found: %s", idStr), http.StatusNotFound)
			return
		}

		// Create new silence
		if p.ByName("id") == "new" {
			m, err := revere.LoadMonitors(db)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable load new silence page: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			err = executeTemplate(w, "silences-edit.html", silenceDataWith(map[string]interface{}{
				"Monitors":    m,
				"Breadcrumbs": silencesIndexBcs(),
			}))

			if err != nil {
				http.Error(w, fmt.Sprintf("Unable load new silence page: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}

		// Edit existing silence
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Silence not found: %s", idStr), http.StatusNotFound)
			return
		}

		s, err := revere.LoadSilence(db, uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable load edit silence page: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		err = executeTemplate(w, "silences-edit.html",
			silenceDataWith(map[string]interface{}{
				"Silence":     s,
				"Breadcrumbs": silencesViewBcs(s.Id, s.MonitorName),
			}))
	}
}

func SilencesSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := getSilenceId(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load silence: %s", err.Error()),
				http.StatusNotFound)
			return
		}

		var s *revere.Silence
		err = json.NewDecoder(req.Body).Decode(&s)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		s.Id = id

		errs, err := s.Validate(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if len(errs) > 0 {
			writeJsonResponse(w, "save silence", map[string]interface{}{"errors": errs})
			return
		}

		err = s.Save(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		writeJsonResponse(w, "save silence", map[string]interface{}{"id": s.Id})
	}
}

func getSilenceId(idStr string) (uint, error) {
	if idStr == "new" {
		return 0, nil
	}

	id, err := strconv.Atoi(idStr)
	return uint(id), err
}

func silenceDataWith(d map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"Title": "silences",
	}
	for k, v := range d {
		data[k] = v
	}
	return data
}
