package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"

	"github.com/yext/revere/datasource"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"
)

func DataSourcesIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		viewmodels, err := datasource.All(DB)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve data sources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		saveStatus, err := getFlash(w, req, "saveStatus")
		if err != nil {
			log.Errorf("Unable to load flash cookie for data sources: %s", err.Error())
		}

		renderable := renderables.NewDataSourceIndex(viewmodels, saveStatus)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve data sources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func DataSourcesSave(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var dss []*datasource.VM
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data sources must be in correct format: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body.Bytes(), &dss)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data sources must be in correct format: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		var errs []string
		for _, ds := range dss {
			errs = append(errs, ds.Validate()...)
		}
		if errs != nil {
			errors, err := json.Marshal(map[string][]string{"errors": errs})
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to save data sources: %s", err.Error()),
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
				http.Error(w, "Unable to retrieve monitors to check data source usage", http.StatusInternalServerError)
			}
			for _, ds := range dss {
				isUsed := isDataSourceInUse(ds.SourceID, monitors, tx)
				if isUsed && ds.Delete {
					http.Error(w, fmt.Sprintf("Can't delete a data source currently used by a monitor. ID: %d", ds.SourceID),
						http.StatusBadRequest)
					return nil
				}
				err = ds.Save(tx)
				if err != nil {
					http.Error(w, fmt.Sprintf("Unable to save data sources: %s", err.Error()),
						http.StatusInternalServerError)
					return nil
				}
			}
			return nil
		})
		dsi := make([]vm.NamedComponent, len(dss))
		for i, ds := range dss {
			dsi[i] = ds
		}
		logSaveArray(dsi, body.Bytes(), req.URL.String())

		setFlash(w, "saveStatus", []byte("updated"))
	}
}

func LoadValidDataSources(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		pt, err := strconv.Atoi(p.ByName("probeType"))

		probe, err := probes.Blank(db.ProbeType(pt))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load probe: %s", err.Error()),
				http.StatusNotFound)
			return
		}

		acceptedTypes := probe.AcceptedSourceTypes()
		sources, err := datasource.AllOfTypes(DB, acceptedTypes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to load data sources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(sources)
		return
	}
}

func isDataSourceInUse(id db.DatasourceID, monitors []*vm.Monitor, tx *db.Tx) bool {
	for _, monitor := range monitors {
		if monitor.Probe.HasDatasource(id) {
			return true
		}
	}
	return false
}
