package server

import (
	"net/http"
	"strconv"

	"github.com/yext/revere/env"
	"github.com/yext/revere/web"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type WebServer struct {
	*env.Env
	router *httprouter.Router
}

func New(env *env.Env) *WebServer {
	router := httprouter.New()
	router.GET("/", web.ActiveIssues(env.DB))
	router.GET("/datasources", web.DataSourcesIndex(env.DB))
	router.POST("/datasources", web.DataSourcesSave(env.DB))
	router.GET("/monitors", web.MonitorsIndex(env.DB))
	router.GET("/monitors/:id", web.MonitorsView(env.DB))
	router.GET("/monitors/:id/edit", web.MonitorsEdit(env.DB))
	router.POST("/monitors/:id/edit", web.MonitorsSave(env.DB))
	router.GET("/monitors/:id/subprobes", web.SubprobesIndex(env.DB))
	router.GET("/monitors/:id/subprobes/:subprobeId", web.SubprobesView(env.DB))
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
	router.ServeFiles("/static/css/*filepath", http.Dir("web/css"))
	router.ServeFiles("/static/js/*filepath", http.Dir("web/js"))
	router.HandlerFunc("GET", "/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/favicon.ico")
	})

	return &WebServer{
		Env:    env,
		router: router,
	}
}

func (w *WebServer) run() {
	port := strconv.Itoa(w.Port)
	log.Info("Listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, w.router))
}

func (w *WebServer) Start() {
	go w.run()
}

func (w *WebServer) Stop() {
	// TODO(fchen): implement graceful shutdown
}
