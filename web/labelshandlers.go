package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/yext/revere"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"
)

func LabelsIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		labels, err := vm.AllLabels(db)
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

func LabelsView(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")

		if id == "new" {
			http.Redirect(w, req, "/labels/new/edit", http.StatusMovedPermanently)
			return
		}

		viewmodel, err := loadLabelViewModel(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewLabelView(viewmodel)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func LabelsEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Label not found", http.StatusNotFound)
			return
		}

		viewmodel, err := loadLabelViewModel(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		monitors, err := vm.AllMonitors(db)
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

func LabelsSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var l *revere.Label
		err := json.NewDecoder(req.Body).Decode(&l)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		errs := l.Validate(db)
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
		err = l.Save(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		redirect, err := json.Marshal(map[string]string{"redirect": fmt.Sprintf("/labels/%d", l.LabelId)})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save label: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(redirect)
	}
}

func loadLabelViewModel(db *sql.DB, unparsedId string) (*vm.Label, error) {
	if unparsedId == "new" {
		viewmodel, err := vm.BlankLabel(db)
		if err != nil {
			return nil, err
		}
		return viewmodel, nil
	}

	id, err := strconv.Atoi(unparsedId)
	if err != nil {
		return nil, err
	}

	viewmodel, err := vm.NewLabel(db, revere.LabelID(id))
	if err != nil {
		return nil, err
	}

	return viewmodel, nil
}
