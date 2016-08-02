package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"
)

func DataSourcesIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		viewmodels, err := datasources.All(DB)
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
		var dss []datasources.VM
		err := json.NewDecoder(req.Body).Decode(&dss)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data sources must be in correct format: %s", err), http.StatusInternalServerError)
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

		err = DB.Tx(func(tx *db.Tx) error {
			monitors, err := vm.AllMonitors(tx)
			if err != nil {
				return err
			}
			for _, ds := range dss {
				isUsed := isDataSourceInUse(ds.SourceID, monitors, tx)
				if isUsed && ds.Delete {
					return errors.Errorf("Can't delete a data source currently used by a monitor. ID: %d", ds.SourceID)
				}
				err = ds.Save(tx)
				if err != nil {
					return errors.Trace(err)
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save data sources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

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
		sources, err := datasources.AllOfTypes(DB, acceptedTypes)
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
