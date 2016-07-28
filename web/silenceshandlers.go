package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"

	"github.com/yext/revere/db"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"
)

func SilencesIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var silences []*vm.Silence
		err := DB.Tx(func(tx *db.Tx) error {
			var err error
			silences, err = vm.AllSilences(tx)
			return errors.Trace(err)
		})
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

func SilencesView(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "new" {
			http.Redirect(w, req, "/silences/new/edit", http.StatusMovedPermanently)
			return
		}

		silence, err := loadSilenceViewModel(DB, id)
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

func SilencesEdit(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			http.Error(w, "Silence not found", http.StatusNotFound)
			return
		}
		query := req.URL.Query()

		silence, err := loadSilenceViewModel(DB, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = silence.SetHtmlParams(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to set parameters from query string: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		// TODO(fchen): consider making single silence load take a tx, or dbOrTx?

		var allMonitors []*vm.Monitor
		err = DB.Tx(func(tx *db.Tx) error {
			var err error
			allMonitors, err = vm.AllMonitors(tx)
			return errors.Trace(err)
		})
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

func SilencesSave(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var s *vm.Silence
		err := json.NewDecoder(req.Body).Decode(&s)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		errs := s.Validate(DB)
		if len(errs) > 0 {
			writeJsonResponse(w, "save silence", map[string]interface{}{"errors": errs})
			return
		}

		err = DB.Tx(func(tx *db.Tx) error {
			return s.Save(tx)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		writeJsonResponse(w, "save silence", map[string]interface{}{"id": s.SilenceID})
	}
}

func loadSilenceViewModel(DB *db.DB, unparsedId string) (*vm.Silence, error) {
	if unparsedId == "new" {
		viewmodel := vm.BlankSilence()
		return viewmodel, nil
	}

	id, err := strconv.Atoi(unparsedId)
	if err != nil {
		return nil, errors.Trace(err)
	}

	viewmodel, err := vm.NewSilence(DB, db.SilenceID(id))
	if err != nil {
		return nil, errors.Trace(err)
	}

	return viewmodel, nil
}

func RedirectToSilence(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		params := req.URL.Query()
		monitorID, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to redirect to silence: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		subprobe := params.Get("subprobe")

		v := url.Values{}
		v.Set("monitorId", strconv.Itoa(monitorID))
		v.Set("subprobes", fmt.Sprintf("^%s$", regexp.QuoteMeta(subprobe)))

		sID := vm.LoadActiveSilence(DB, db.MonitorID(monitorID), subprobe)
		var url string
		if sID != 0 {
			url = fmt.Sprintf("/silences/%d/edit?%s", sID, v.Encode())
		} else {
			url = fmt.Sprintf("/silences/new/edit?%s", v.Encode())
		}

		http.Redirect(w, req, url, http.StatusMovedPermanently)
		return
	}
}
