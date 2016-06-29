package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/juju/errors"
	"github.com/yext/revere/datasources"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probes"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func DataSourcesIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		viewmodels, err := datasources.All(DB)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve data sources: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		renderable := renderables.NewDataSourceIndex(viewmodels)
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
			for _, ds := range dss {
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

		http.Redirect(w, req, "/datasources", http.StatusMovedPermanently)
		return
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
