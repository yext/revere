package web

import (
	"database/sql"
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
		err = executeTemplate(w, "silences-index.html",
			map[string]interface{}{
				"Title":       "silences",
				"Silences":    s,
				"Breadcrumbs": silencesIndexBcs(),
			})
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
			map[string]interface{}{
				"Title":       "silences",
				"Silence":     s,
				"Breadcrumbs": silencesViewBcs(s.Id, s.MonitorName),
			})
		if err != nil {
			fmt.Println("Got err executing template:", err.Error())
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func SilencesEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	}
}
