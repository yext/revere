package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yext/revere"
	"github.com/yext/revere/datasources"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func DataSourcesIndex(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		viewmodels, err := loadAllDataSourceTypeViewModels(db)
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

func DataSourcesSave(db *sql.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var dataSources []*revere.DataSource
		err := json.NewDecoder(req.Body).Decode(&dataSources)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data sources must be in correct format: %s", err), http.StatusInternalServerError)
			return
		}

		for _, dataSource := range dataSources {
			err = dataSource.Save(db)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to save data sources: %s", err.Error()),
					http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, req, "/datasources", http.StatusMovedPermanently)
		return
	}
}

func loadAllDataSourceTypeViewModels(db *sql.DB) (models []*vm.DataSourceTypeViewModel, err error) {
	for _, dst := range datasources.GetDataSourceTypes() {
		dstvm, err := vm.NewDataSourceTypeViewModel(db, dst)
		if err != nil {
			return nil, err
		}
		models = append(models, dstvm)
	}
	return
}
