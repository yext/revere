package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/julienschmidt/httprouter"

	"github.com/yext/revere/db"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/resource"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"
)

func ResourcesIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		viewmodels, err := resource.All(DB)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve resources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		saveStatus, err := getFlash(w, req, "saveStatus")
		if err != nil {
			log.Errorf("Unable to load flash cookie for resources: %s", err.Error())
		}

		renderable := renderables.NewResourcesIndex(viewmodels, saveStatus)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve resources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func ResourcesSave(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var rs []*resource.VM
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Resources must be in correct format: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body.Bytes(), &rs)
		if err != nil {
			http.Error(w, fmt.Sprintf("Resources must be in correct format: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		var errs []string
		for _, r := range rs {
			errs = append(errs, r.Validate()...)
		}
		if errs != nil {
			errors, err := json.Marshal(map[string][]string{"errors": errs})
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to save resources: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(errors)
			return
		}

		DB.Tx(func(tx *db.Tx) error {
			monitors, err := vm.AllMonitors(tx)
			if err != nil {
				http.Error(w, "Unable to retrieve monitors to check resource usage", http.StatusInternalServerError)
			}
			for _, r := range rs {
				isUsed := isResourceInUse(r.ResourceID, monitors, tx)
				if isUsed && r.Delete {
					http.Error(w, fmt.Sprintf("Can't delete a resource currently used by a monitor. ID: %d", r.ResourceID),
						http.StatusBadRequest)
					return nil
				}
				err = r.Save(tx)
				if err != nil {
					http.Error(w, fmt.Sprintf("Unable to save resources: %s", err.Error()),
						http.StatusInternalServerError)
					return nil
				}
			}
			return nil
		})
		ri := make([]vm.NamedComponent, len(rs))
		for i, r := range rs {
			ri[i] = r
		}
		logSaveArray(ri, body.Bytes(), req.URL.String())

		setFlash(w, "saveStatus", []byte("updated"))
	}
}

func LoadValidResources(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		pt, err := strconv.Atoi(p.ByName("probeType"))

		blankProbe, err := probe.Blank(db.ProbeType(pt))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load probe: %s", err.Error()),
				http.StatusNotFound)
			return
		}

		acceptedTypes := blankProbe.AcceptedResourceTypes()
		sources, err := resource.AllOfTypes(DB, acceptedTypes)

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load resources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(sources)
		return
	}
}

func isResourceInUse(id db.ResourceID, monitors []*vm.Monitor, tx *db.Tx) bool {
	for _, monitor := range monitors {
		if monitor.Probe.HasResource(id) {
			return true
		}
	}
	return false
}
