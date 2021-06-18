// Package server implements the core engine for Revere's web mode.
package server

import (
	"net/http"
	"strconv"

	"github.com/braintree/manners"
	"github.com/yext/revere/boxes"
	"github.com/yext/revere/env"
	"github.com/yext/revere/web"

	log "github.com/sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

// WebServer wraps a Router + DB object, and allows users to configure Revere
// components through a web UI.
type WebServer struct {
	*env.Env
	router  *httprouter.Router
	stopped chan struct{}
}

// New initializes the WebServer.
func New(env *env.Env) *WebServer {
	cssFiles := boxes.CSS()
	jsFiles := boxes.JS()
	favicon := boxes.Favicon()

	router := httprouter.New()
	router.GET("/", web.ActiveIssues(env.DB))
	router.GET("/resources", web.ResourcesIndex(env.DB))
	router.GET("/resources/probe/:probeType", web.LoadValidResources(env.DB))
	router.POST("/resources", web.ResourcesSave(env.DB))
	router.GET("/resourcetype/:id", web.LoadResourceTemplate(env.DB))
	router.GET("/monitors", web.MonitorsIndex(env.DB))
	router.GET("/monitors/:id", web.MonitorsView(env.DB))
	router.GET("/monitors/:id/edit", web.MonitorsEdit(env.DB))
	router.POST("/monitors/:id/edit", web.MonitorsSave(env.DB))
	router.GET("/monitors/:id/subprobes", web.SubprobesIndex(env.DB))
	router.GET("/monitors/:id/subprobes/:subprobeId", web.SubprobesView(env.DB))
	router.DELETE("/monitors/:id/subprobes/:subprobeId/delete", web.DeleteSubprobe(env.DB));
	router.GET("/monitors/:id/probe/edit/:probeType", web.LoadProbeTemplate(env.DB))
	router.GET("/monitors/:id/target/edit/:targetType", web.LoadTargetTemplate)
	router.GET("/silences", web.SilencesIndex(env.DB))
	router.GET("/silences/:id", web.SilencesView(env.DB))
	router.GET("/silences/:id/edit", web.SilencesEdit(env.DB))
	router.POST("/silences/:id/edit", web.SilencesSave(env.DB))
	router.GET("/labels", web.LabelsIndex(env.DB))
	router.GET("/labels/:id", web.LabelsView(env.DB))
	router.GET("/labels/:id/edit", web.LabelsEdit(env.DB))
	router.POST("/labels/:id/edit", web.LabelsSave(env.DB))
	router.GET("/settings", web.SettingsIndex(env.DB))
	router.POST("/settings", web.SettingsSave(env.DB))
	router.GET("/redirectToSilence", web.RedirectToSilence(env.DB))

	router.ServeFiles("/static/css/*filepath", cssFiles.HTTPBox())
	router.ServeFiles("/static/js/*filepath", jsFiles.HTTPBox())
	router.Handler("GET", "/favicon.ico", http.FileServer(favicon.HTTPBox()))

	return &WebServer{
		Env:     env,
		router:  router,
		stopped: make(chan struct{}),
	}
}

func (w *WebServer) run() {
	port := strconv.Itoa(int(w.Port))
	log.Info("Listening on :", port)
	err := manners.ListenAndServe(":"+port, w.router)
	if err != nil {
		log.Fatal(err)
	}
	close(w.stopped)
}

// Start starts the WebServer.
func (w *WebServer) Start() {
	go w.run()
}

// Stop gracefully stops the WebServer. Stop will block until the WebServer
// finishes handling its current requests, and will then shut down the
// WebServer.
func (w *WebServer) Stop() {
	manners.Close()
	<-w.stopped
}
