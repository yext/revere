package web

// TODO(dp): rename this file once we finish migration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/targets"
	"github.com/yext/revere/util"

	"github.com/julienschmidt/httprouter"
)

var (
	tMap    map[string]*template.Template = make(map[string]*template.Template)
	funcMap template.FuncMap              = make(template.FuncMap)
)

func init() {
	funcMap["isLastBc"] = isLastBc
	funcMap["strEq"] = util.StrEq
	funcMap["targets"] = targets.AllTargets
	funcMap["targetScripts"] = targetScripts
	funcMap["probes"] = probes.AllProbes
	funcMap["probeScripts"] = probeScripts
}

func LoadTemplates() {
	templateInfo, err := ioutil.ReadDir("web/views")
	for _, t := range templateInfo {
		if t.IsDir() || !strings.HasSuffix(t.Name(), ".html") {
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
				"Title":       "Active Issues",
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
