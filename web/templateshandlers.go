package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/targets"
	"github.com/yext/revere/web/vm/renderables"
)

func LoadDataSourceTemplate(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		i, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, "Id must be an int", http.StatusInternalServerError)
			return
		}
		id := db.SourceType(i)
		ds, err := datasources.Blank(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("No data source type with id: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		dsvm := &datasources.VM{
			SourceType: id,
			DataSource: ds,
		}
		dsv := renderables.NewDataSourceView(dsvm)

		tmpl, err := renderables.RenderPartial(dsv)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to render template: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		template, err := json.Marshal(map[string]template.HTML{"template": tmpl})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load data source: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(template)
	}
}

func LoadProbeTemplate(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		pt, err := strconv.Atoi(p.ByName("probeType"))

		if err != nil {
			http.Error(w, fmt.Sprintf("Probe type not found: %s", p.ByName("probeType")), http.StatusNotFound)
			return
		}

		probe, err := probes.Blank(db.ProbeType(pt))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load probe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		pe := renderables.NewProbeEdit(probe)

		tmpl, err := renderables.RenderPartial(pe)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load probe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		template, err := json.Marshal(map[string]template.HTML{"template": tmpl})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load probe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(template)
	}
}

func LoadTargetTemplate(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	tt, err := strconv.Atoi(p.ByName("targetType"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Target type not found: %s", p.ByName("targetType")), http.StatusNotFound)
		return
	}

	target, err := targets.Blank(db.TargetType(tt))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load target: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	te := renderables.NewTargetEdit(target)

	tmpl, err := renderables.RenderPartial(te)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load target: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	template, err := json.Marshal(map[string]template.HTML{"template": tmpl})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load target: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(template)
}
