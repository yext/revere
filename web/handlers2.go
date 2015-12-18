package web

// TODO(dp): rename this file once we finish migration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/yext/revere"
	"github.com/yext/revere/util"

	"github.com/julienschmidt/httprouter"
)

const (
	format = "01/02/2006 3:04 PM"
)

var (
	tMap    map[string]*template.Template = make(map[string]*template.Template)
	funcMap template.FuncMap              = make(template.FuncMap)
)

func init() {
	funcMap["isLastBc"] = isLastBc
	funcMap["strEq"] = util.StrEq
}

func LoadTemplates() {
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
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
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

func writeJsonResponse(w http.ResponseWriter, action string, data map[string]interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to %s: %s", action, err.Error()),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func executeTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	if _, ok := data["States"]; !ok {
		data["States"] = revere.ReverseStates
	}
	return tMap[name].ExecuteTemplate(w, name, data)
}
