package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func SubprobesIndex(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", p.ByName("id")),
				http.StatusNotFound)
			return
		}

		var (
			subprobes []*vm.Subprobe
			monitor   *vm.Monitor
		)
		err = DB.Tx(func(tx *db.Tx) error {
			var err error
			subprobes, err = vm.AllSubprobesFromMonitor(tx, db.MonitorID(id))
			if err != nil {
				return errors.Trace(err)
			}
			monitor, err = vm.NewMonitor(tx, db.MonitorID(id))
			return errors.Trace(err)
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewSubprobesIndex(subprobes, monitor)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobes: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func SubprobesView(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		mId, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Monitor not found: %s", p.ByName("id")),
				http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(p.ByName("subprobeId"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %s", p.ByName("subprobeId")),
				http.StatusNotFound)
			return
		}

		subprobe, err := vm.NewSubprobe(DB, db.SubprobeID(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		if subprobe.MonitorID != db.MonitorID(mId) {
			http.Error(w, fmt.Sprintf("Subprobe %d does not exist for monitor: %d", id, mId),
				http.StatusNotFound)
			return
		}

		probe, err := vm.NewProbe(DB, subprobe.MonitorID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to get probe for monitor: %s", err.Error()),
				http.StatusNotFound)
		}

		readings, err := vm.AllReadingsFromSubprobe(DB, db.SubprobeID(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve readings: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewSubprobeView(probe, subprobe, readings)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve subprobe: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
	}
}

func DeleteSubprobe(DB *db.DB) func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		subprobeId, err := strconv.Atoi(p.ByName("subprobeId"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Subprobe not found: %s", p.ByName("subprobeId")),
				http.StatusNotFound);
		}
		err = DB.Tx(func(tx *db.Tx) error {
			var err error
			err = vm.DeleteSubprobe(tx, subprobeId)
			if err != nil {
				return errors.Trace(err)
			}
			return nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to delete subprobes: %v", err),
				http.StatusInternalServerError)
			return
		}
	}
}