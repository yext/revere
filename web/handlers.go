/*
Package web handles all web-based user interaction. Users can view and edit
Revere components and configurations through its web server. A typical request
is routed through the server package to a handler method in the web package.
Handler methods then create and manipulate the appropriate VMs, and use the
renderables package to display the respective web page.
*/
package web

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/probe"
	"github.com/yext/revere/setting"
	"github.com/yext/revere/state"
	"github.com/yext/revere/target"
	"github.com/yext/revere/web/tmpl"
	"github.com/yext/revere/web/vm"
	"github.com/yext/revere/web/vm/renderables"

	"github.com/julienschmidt/httprouter"
)

func init() {
	tmpl.AddDefaultFunc("isLastBc", vm.IsLastBc)
	tmpl.AddDefaultFunc("setTitle", tmpl.SetTitle)
	tmpl.AddDefaultFunc("strEq", tmpl.StrEq)
	tmpl.AddDefaultFunc("deepEq", reflect.DeepEqual)
	tmpl.AddDefaultFunc("stateClass", state.CSSClass)
	tmpl.AddDefaultFunc("targets", target.AllTargets)
	tmpl.AddDefaultFunc("settings", setting.AllTypes)
	tmpl.AddDefaultFunc("probeTypes", probe.AllTypes)
}

func ActiveIssues(DB *db.DB) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var (
			err           error
			subprobes     []*vm.Subprobe
			monitorLabels map[db.MonitorID][]*vm.MonitorLabel
			allLabels     []*vm.Label
		)

		l := req.FormValue("label")
		labelID, err := strconv.Atoi(l)
		labelUsed := err == nil

		err = DB.Tx(func(tx *db.Tx) (err error) {
			if labelUsed {
				subprobes, err = vm.AllAbnormalSubprobesForLabel(tx, db.LabelID(labelID))
			} else {
				subprobes, err = vm.AllAbnormalSubprobes(tx)
			}
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
					http.StatusInternalServerError)
				return errors.Trace(err)
			}

			monitorLabels, err = vm.AllMonitorLabelsForSubprobes(tx, subprobes)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
					http.StatusInternalServerError)
				return errors.Trace(err)
			}

			allLabels, err = vm.AllLabels(tx)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to retrieve labels: %s", err.Error()),
					http.StatusInternalServerError)
				return errors.Trace(err)
			}
			return nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		renderable := renderables.NewActiveIssues(subprobes, allLabels, monitorLabels)
		err = render(w, renderable)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to retrieve active issues: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		return
	}
}

func writeJsonResponse(w http.ResponseWriter, action string, data map[string]interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to %s: %s", action, err.Error()),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func render(w io.Writer, r renderables.Renderable) error {
	return renderables.Render(w, r)
}

func setFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Name: name, Value: encode(value), Path: "/"}
	http.SetCookie(w, c)
}

func getFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return nil, err
	}
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0), Path: "/"}
	http.SetCookie(w, dc)
	return value, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func logSave(c vm.NamedComponent, body []byte, url string) {
	b := new(bytes.Buffer)
	err := json.Indent(b, body, "", "\t")
	if err != nil {
		log.WithFields(log.Fields{
			"json":      string(body),
			"URL":       url,
			"component": c,
		}).Error("Unable to indent request for logging")
		return
	}
	log.WithFields(log.Fields{
		"Component": c.ComponentName(),
		"ID":        c.Id(),
		"URL":       url,
	}).Info(b.String())
}

func logSaveArray(c []vm.NamedComponent, body []byte, url string) {
	b := new(bytes.Buffer)
	if len(c) > 0 {
		first := c[0]
		err := json.Indent(b, body, "", "\t")
		if err != nil {
			log.WithFields(log.Fields{
				"json":       string(body),
				"URL":        url,
				"components": c,
			}).Error("Unable to indent request for logging")
			return
		}
		log.WithFields(log.Fields{
			"Component": first.ComponentName(),
			"URL":       url,
		}).Info(b.String())
	}
}
