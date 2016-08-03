package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/setting"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func SettingsIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		viewmodels, err := setting.All(DB)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve settings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		saveStatus, err := getFlash(w, req, "saveStatus")
		if err != nil {
			log.Errorf("Unable to load flash cookie for settings: %s", err.Error())
		}

		renderable := renderables.NewSettingsIndex(viewmodels, saveStatus)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve settings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SettingsSave(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var ss []*setting.VM
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save settings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body.Bytes(), &ss)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save seetings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		var errs []string
		for _, s := range ss {
			errs = append(errs, s.Validate()...)
		}
		if errs != nil {
			errors, err := json.Marshal(map[string][]string{"errors": errs})
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to save settings: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(errors)
			return
		}

		err = DB.Tx(func(tx *db.Tx) error {
			for _, s := range ss {
				err := s.Save(tx)
				if err != nil {
					return errors.Trace(err)
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save settings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		si := make([]vm.NamedComponent, len(ss))
		for i, s := range ss {
			si[i] = s
		}
		logSaveArray(si, body.Bytes(), req.URL.String())

		setFlash(w, "saveStatus", []byte("updated"))
	}
}
