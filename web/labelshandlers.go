package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/yext/revere/db"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"
)

func LabelsIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var labels []*vm.Label
		err := DB.Tx(func(tx *db.Tx) error {
			var err error
			labels, err = vm.AllLabels(tx)
			return errors.Trace(err)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewLabelsIndex(labels)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func LabelsView(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")

		if id == "new" {
			http.Redirect(w, req, "/labels/new/edit", http.StatusMovedPermanently)
			return
		}

		viewmodel, err := loadLabelViewModel(DB, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		saveStatus, err := getFlash(w, req, "saveStatus")
		if err != nil {
			log.Errorf("Unable to load flash cookie for label: %s", err.Error())
		}

		renderable := renderables.NewLabelView(viewmodel, saveStatus)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func LabelsEdit(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Label not found", http.StatusNotFound)
			return
		}

		viewmodel, err := loadLabelViewModel(DB, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		var monitors []*vm.Monitor
		err = DB.Tx(func(tx *db.Tx) error {
			var err error
			monitors, err = vm.AllMonitors(tx)
			return errors.Trace(err)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve monitors for label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewLabelEdit(viewmodel, monitors)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func LabelsSave(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var l *vm.Label
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body.Bytes(), &l)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		errs := l.Validate(DB)
		if errs != nil {
			errors, err := json.Marshal(map[string][]string{"errors": errs})
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(errors)
			return
		}

		var saveStatus string
		if l.IsCreate() {
			saveStatus = "created"
		} else {
			saveStatus = "updated"
		}

		err = DB.Tx(func(tx *db.Tx) error {
			err := l.Save(tx)
			return errors.Trace(err)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		logSave(l, body.Bytes(), req.URL.String())

		redirect, err := json.Marshal(map[string]string{"redirect": fmt.Sprintf("/labels/%d", l.LabelID)})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		setFlash(w, "saveStatus", []byte(saveStatus))

		w.Header().Set("Content-Type", "application/json")
		w.Write(redirect)
	}
}

func loadLabelViewModel(DB *db.DB, unparsedId string) (*vm.Label, error) {
	var viewmodel *vm.Label
	if unparsedId == "new" {
		viewmodel = vm.BlankLabel()
		return viewmodel, nil
	}

	id, err := strconv.Atoi(unparsedId)
	if err != nil {
		return nil, errors.Trace(err)
	}

	err = DB.Tx(func(tx *db.Tx) error {
		var err error
		viewmodel, err = vm.NewLabel(tx, db.LabelID(id))
		return errors.Trace(err)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	return viewmodel, nil
}
