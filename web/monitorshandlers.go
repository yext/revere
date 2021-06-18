package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
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
				return errors.Trace(err)
			}

			err = vm.PopulateLabelsForMonitors(tx, monitors)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
					http.StatusInternalServerError)
				return errors.Trace(err)
			}

			labels, err = vm.AllLabels(tx)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
					http.StatusInternalServerError)
				return errors.Trace(err)
			}

			return nil
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
			monitor, err = loadMonitorViewModel(tx, id)
			return errors.Trace(err)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		saveStatus, err := getFlash(w, req, "saveStatus")
		if err != nil {
			log.Errorf("Unable to load flash cookie for monitor: %s", err.Error())
		}

		renderable := renderables.NewMonitorView(monitor, saveStatus)
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

			return nil
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
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body.Bytes(), &m)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		errs := m.Validate(DB)
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

		var saveStatus string
		if m.IsCreate() {
			saveStatus = "created"
		} else {
			saveStatus = "updated"
		}

		err = DB.Tx(func(tx *db.Tx) error {
			return m.Save(tx)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		logSave(m, body.Bytes(), req.URL.String())

		redirect, err := json.Marshal(map[string]string{"redirect": fmt.Sprintf("/monitors/%d", m.MonitorID)})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save monitor: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		setFlash(w, "saveStatus", []byte(saveStatus))

		w.Header().Set("Content-Type", "application/json")
		w.Write(redirect)
	}
}

func loadMonitorViewModel(tx *db.Tx, unparsedId string) (*vm.Monitor, error) {
	if unparsedId == "new" {
		blankMonitor, err := vm.BlankMonitor()
		if err != nil {
			return nil, errors.Trace(err)
		}
		return blankMonitor, nil
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
