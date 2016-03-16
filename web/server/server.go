package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/yext/revere"
	"github.com/yext/revere/web"

	"github.com/julienschmidt/httprouter"
)

func main() {
	flag.Parse()
	env, err := revere.BuildEnvFromFile(flag.Arg(0))
	if err != nil {
		fmt.Println("Unable to start web server.")
		return
	}

	db := env.Db()

	router := httprouter.New()

	router.GET("/", web.ActiveIssues(db))
	router.GET("/datasources", web.DataSourcesIndex(db))
	router.POST("/datasources", web.DataSourcesUpdate(db))
	router.POST("/datasources/delete", web.DataSourcesDelete(db))
	router.GET("/monitors", web.MonitorsIndex(db))
	router.GET("/monitors/:id", web.MonitorsView(db))
	router.GET("/monitors/:id/edit", web.MonitorsEdit(db))
	router.POST("/monitors/:id/edit", web.MonitorsSave(db))
	router.GET("/monitors/:id/subprobes", web.SubprobesIndex(db))
	router.GET("/monitors/:id/subprobes/:subprobeId", web.SubprobesView(db))
	router.GET("/monitors/:id/probe/edit/:probeType", web.LoadProbeTemplate(db))
	router.GET("/monitors/:id/target/edit/:targetType", web.LoadTargetTemplate)
	router.GET("/silences", web.SilencesIndex(db))
	router.GET("/silences/:id", web.SilencesView(db))
	router.GET("/silences/:id/edit", web.SilencesEdit(db))
	router.POST("/silences/:id/edit", web.SilencesSave(db))
	router.GET("/labels", web.LabelsIndex(db))
	router.GET("/labels/:id", web.LabelsView(db))
	router.GET("/labels/:id/edit", web.LabelsEdit(db))
	router.POST("/labels/:id/edit", web.LabelsSave(db))
	router.GET("/settings", web.SettingsIndex)
	router.POST("/settings/:id/save", web.SaveSetting)
	router.ServeFiles("/static/css/*filepath", http.Dir("web/css"))
	router.ServeFiles("/static/js/*filepath", http.Dir("web/js"))
	router.HandlerFunc("GET", "/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/favicon.ico")
	})

	web.LoadTemplates()

	port := strconv.Itoa(env.Port())
	fmt.Printf("Listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
