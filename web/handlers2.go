package web

// TODO(dp): rename this file once we finish migration

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

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

// Temporary function. We'll probably need to pass in the db here, so we'll return a function
func MonitorsView(_ *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		http.Error(w, fmt.Sprintf("Not yet implemented for id %s", p.ByName("id")), http.StatusNotImplemented)
	}
}
