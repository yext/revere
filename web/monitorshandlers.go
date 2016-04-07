package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/yext/revere"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func MonitorsIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var (
			monitors []*vm.Monitor
			err      error
		)

		l := req.FormValue("label")
		labelId, err := strconv.Atoi(l)
		if err != nil {
			monitors, err = vm.AllMonitors(db)
		} else {
			monitors, err = vm.AllMonitorsForLabel(db, labelId)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitors: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		err = vm.PopulateLabelsForMonitors(db, monitors)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		labels, err := vm.AllLabels(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewMonitorsIndex(monitors, labels)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitors: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func MonitorsView(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")

		if id == "new" {
			http.Redirect(w, req, "/monitors/new/edit", http.StatusMovedPermanently)
			return
		}

		viewmodel, err := loadMonitorViewModel(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewMonitorView(viewmodel)
		err = render(w, renderable)

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
			http.Error(w, "Monitor not found", http.StatusNotFound)
			return
		}

		viewmodel, err := loadMonitorViewModel(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewMonitorEdit(viewmodel)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
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
		errs := m.Validate(db)
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
		err = m.Save(db)
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

func loadMonitorViewModel(db *sql.DB, unparsedId string) (*vm.Monitor, error) {
	if unparsedId == "new" {
		viewmodel, err := vm.BlankMonitor(db)
		if err != nil {
			return nil, err
		}
		return viewmodel, nil
	}

	id, err := strconv.Atoi(unparsedId)
	if err != nil {
		return nil, err
	}

	viewmodel, err := vm.NewMonitor(db, id)
	if err != nil {
		return nil, err
	}

	return viewmodel, nil
}

func LoadProbeTemplate(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		pt, err := strconv.Atoi(p.ByName("probeType"))

		if err != nil {
			http.Error(w, fmt.Sprintf("Probe type not found: %s", p.ByName("probeType")), http.StatusNotFound)
			return
		}

		probe, err := vm.BlankProbe(db, pt)
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

	target, err := vm.BlankTarget(tt)
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
