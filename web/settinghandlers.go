package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/yext/revere/settings"

	"github.com/julienschmidt/httprouter"
)

func SettingsIndex(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	allSettings := settings.GetAllLoadedSettings()
	templates := make([]template.HTML, len(allSettings))
	for i, setting := range allSettings {
		rendered, err := setting.Render()
		if err != nil {
			rendered = template.HTML(fmt.Sprintf("<p>Failed to render %T : %v </p>", setting, err))
		}
		templates[i] = rendered
	}
	data := map[string]interface{}{
		"Title":    "Settings",
		"Settings": templates,
	}
	err := executeTemplate(w, "settings-index.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading settings : %s", err.Error()),
			http.StatusInternalServerError)
		return
	}
}

func SaveSetting(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid setting id provided: %s", p.ByName("id")),
			http.StatusBadRequest)
		return
	}
	setting, found := settings.GetSetting(id)
	if !found {
		http.Error(w, fmt.Sprintf("Setting not found: %s", p.ByName("id")),
			http.StatusNotFound)
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to save setting: %s", err.Error()),
			http.StatusBadRequest)
		return
	}
	err = setting.Save(string(b))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to save setting: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, "save setting", map[string]interface{}{"id": id})
}
