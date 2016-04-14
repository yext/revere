package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/yext/revere"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/settings"
	"github.com/yext/revere/targets"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

var partials = "web/views/partials/*.html"

func init() {
	tmpl.AddDefaultFunc("isLastBc", vm.IsLastBc)
	tmpl.AddDefaultFunc("setTitle", tmpl.SetTitle)
	tmpl.AddDefaultFunc("strEq", tmpl.StrEq)
	tmpl.AddDefaultFunc("targets", targets.AllTargets)
	tmpl.AddDefaultFunc("probes", probes.AllProbes)
	tmpl.AddDefaultFunc("settings", settings.AllSettingTypes)
	tmpl.SetPartialsLocation(partials)
}

func ActiveIssues(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var (
			subprobes []*vm.Subprobe
			err       error
		)

		l := req.FormValue("label")
		labelId, err := strconv.Atoi(l)
		if err != nil {
			subprobes, err = vm.AllAbnormalSubprobes(db)
		} else {
			subprobes, err = vm.AllAbnormalSubprobesForLabel(db, revere.LabelID(labelId))
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		monitorLabels, err := vm.AllMonitorLabelsForSubprobes(db, subprobes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		allLabels, err := vm.AllLabels(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewActiveIssues(subprobes, allLabels, monitorLabels)
		err = render(w, renderable)
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

func render(w io.Writer, r renderables.Renderable) error {
	return renderables.Render(w, r)
}
