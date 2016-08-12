package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/resource"
	"github.com/yext/revere/target"
	"github.com/yext/revere/web/vm/renderables"
)

func LoadResourceTemplate(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		i, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, "Id must be an int", http.StatusInternalServerError)
			return
		}
		id := db.ResourceType(i)
		ds, err := resource.Blank(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("No resource type with id: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		rvm := &resource.VM{
			ResourceType: id,
			Resource:     ds,
		}
		rv := renderables.NewResourceView(rvm)

		tmpl, err := renderables.RenderPartial(rv)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to render template: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		template, err := json.Marshal(map[string]template.HTML{"template": tmpl})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load resource: %s", err.Error()),
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

		blankProbe, err := probe.Blank(db.ProbeType(pt))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load probe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		pe := renderables.NewProbeEdit(blankProbe)

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

	target, err := target.Blank(db.TargetType(tt))
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
