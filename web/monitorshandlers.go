package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func MonitorsIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		labelId, err := strconv.Atoi(req.FormValue("label"))
		labelUsed := err == nil

		var (
			monitors []*vm.Monitor
			labels   []*vm.Label
		)
		err = DB.Tx(func(tx *db.Tx) error {
			if labelUsed {
				monitors, err = vm.AllMonitorsForLabel(tx, db.LabelID(labelId))
			} else {
				monitors, err = vm.AllMonitors(tx)
			}
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve monitors: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			err = vm.PopulateLabelsForMonitors(tx, monitors)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			labels, err = vm.AllLabels(tx)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitors: %s", err.Error()),
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

func MonitorsView(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")

		if id == "new" {
			http.Redirect(w, req, "/monitors/new/edit", http.StatusMovedPermanently)
			return
		}

		var monitor *vm.Monitor
		err := DB.Tx(func(tx *db.Tx) (err error) {
			monitor, err = loadMonitorViewModel(DB, id)
			return errors.Trace(err)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewMonitorView(monitor)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func MonitorsEdit(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Monitor not found", http.StatusNotFound)
			return
		}

		var (
			monitor *vm.Monitor
			labels  []*vm.Label
		)
		err := DB.Tx(func(tx *db.Tx) (err error) {
			monitor, err = loadMonitorViewModel(tx, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			labels, err = vm.AllLabels(tx)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve labels for monitor: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewMonitorEdit(monitor, labels)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func MonitorsSave(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var m *vm.Monitor
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

		err = DB.Tx(func(tx *db.Tx) error {
			return m.Save(tx)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		redirect, err := json.Marshal(map[string]string{"redirect": fmt.Sprintf("/monitors/%d", m.MonitorId)})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(redirect)
	}
}

func loadMonitorViewModel(tx *db.Tx, unparsedId string) (*vm.Monitor, error) {
	if unparsedId == "new" {
		return vm.BlankMonitor(), nil
	}

	id, err := strconv.Atoi(unparsedId)
	if err != nil {
		return nil, errors.Trace(err)
	}

	monitor, err := vm.NewMonitor(tx, db.MonitorID(id))
	if err != nil {
		return nil, errors.Trace(err)
	}

	return monitor, nil
}
