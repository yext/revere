package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yext/revere"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func SilencesIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		silences, err := vm.AllSilences(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silences: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewSilencesIndex(silences)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silences: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SilencesView(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "new" {
			http.Redirect(w, req, "/silences/new/edit", http.StatusMovedPermanently)
			return
		}

		silence, err := loadSilenceViewModel(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewSilenceView(silence)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SilencesEdit(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Silence not found", http.StatusNotFound)
			return
		}

		silence, err := loadSilenceViewModel(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		allMonitors, err := vm.AllMonitors(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewSilenceEdit(silence, allMonitors)
		err = render(w, renderable)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SilencesSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var s *revere.Silence
		err := json.NewDecoder(req.Body).Decode(&s)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		oldS, err := revere.LoadSilence(db, s.SilenceId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		errs := s.ValidateAgainstOld(oldS)
		if len(errs) > 0 {
			writeJsonResponse(w, "save silence", map[string]interface{}{"errors": errs})
			return
		}

		err = s.Save(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		writeJsonResponse(w, "save silence", map[string]interface{}{"id": s.SilenceId})
	}
}

func loadSilenceViewModel(db *sql.DB, unparsedId string) (*vm.Silence, error) {
	if unparsedId == "new" {
		viewmodel, err := vm.BlankSilence(db)
		if err != nil {
			return nil, err
		}

		return viewmodel, nil
	}

	id, err := strconv.Atoi(unparsedId)
	if err != nil {
		return nil, err
	}

	viewmodel, err := vm.NewSilence(db, revere.SilenceID(id))
	if err != nil {
		return nil, err
	}

	return viewmodel, nil
}
