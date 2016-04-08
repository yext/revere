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

func SettingsIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

		viewmodels, err := vm.AllSettings(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve settings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		renderable := renderables.NewSettingsIndex(viewmodels)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve settings: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SettingsSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var settings []*revere.Setting
		err := json.NewDecoder(req.Body).Decode(&settings)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save settings: %s", err), http.StatusInternalServerError)
			return
		}

		var errs []string
		for _, setting := range settings {
			errs = append(errs, setting.Validate()...)
			if errs == nil {
				err = setting.Save(db)
				if err != nil {
					http.Error(w, fmt.Sprintf("Unable to save settings: %s", err.Error()),
						http.StatusInternalServerError)
					return
				}
			}
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

		http.Redirect(w, req, "/settings", http.StatusMovedPermanently)
		return
	}
}

func LoadSettingTemplate(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	st, err := strconv.Atoi(p.ByName("settingType"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Setting type not found: %s", p.ByName("settingType")), http.StatusNotFound)
		return
	}

	setting, err := vm.BlankSetting(st)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load setting: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	se := renderables.NewSettingEdit(setting)

	tmpl, err := renderables.RenderPartial(se)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load setting: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	template, err := json.Marshal(map[string]template.HTML{"template": tmpl})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load setting: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(template)
}
