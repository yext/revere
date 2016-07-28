/*
Revere is an alerting system for medium-sized microservices architectures.

Usage:

	revere [-conf env.json] [-mode mode,...]

Configuration

The -conf flag specifies the path to a JSON file configuring Revere's static
environment. When no file is specified, Revere uses its default settings.

See github.com/yext/revere/env.EnvJSONModel for the structure the JSON in this
file should take.

Modes

The -mode flag specifies the comma-separated modes that this invocation of
Revere should run.

The initdb mode initializes Revere's database storage. Depending on whether
there are existing Revere tables in the database specified by the environment
configuration, this mode either creates a new storage area from scratch or
updates an existing area to the current schema. This mode cannot be combined
with any other modes.

The daemon mode runs the daemon that monitors systems and generates alerts.

The web mode serves the HTTP UI for administering Revere and viewing its
current state.

The -mode flag defaults to daemon,web.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
	"golang.org/x/sys/unix"

	"github.com/yext/revere/daemon"
	"github.com/yext/revere/env"
	"github.com/yext/revere/web/server"
)

var (
	conf     = flag.String("conf", "", "JSON `file` configuring Revere's static environment")
	mode     = flag.String("mode", "daemon,web", "comma-separated `modes` to run")
	logLevel = flag.String("logLevel", "warn", "Logrus `level` to log at")
)

const (
	SQL_SETUP_SCRIPT = "conf/setup/initialize_db.sql"
)

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected arguments: %v\n", flag.Args())
		flag.Usage()
		os.Exit(2)
	}

	err := initLog()
	ifErrPrintAndExit(err)

	env, err := loadEnv()
	ifErrPrintAndExit(err)

	modes, err := parseMode()
	ifErrPrintAndExit(err)

	if modes[0] == "initdb" {
		err := initDB(env)
		ifErrPrintAndExit(err)
		return
	}

	for _, mode := range modes {
		switch mode {
		case "daemon":
			d := daemon.New(env)
			d.Start()
			defer d.Stop()
		case "web":
			w := server.New(env)
			w.Start()
			defer w.Stop()
		}
	}

	waitForExitSignal()

	log.Info("Caught exit signal. Revere is shutting down...")
}

func initLog() error {
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		return errors.Mask(err)
	}

	log.SetLevel(level)
	return nil
}

func loadEnv() (*env.Env, error) {
	var json []byte
	var desc string
	if *conf != "" {
		desc = "env conf " + *conf

		var err error
		json, err = ioutil.ReadFile(*conf)
		if err != nil {
			return nil, errors.Maskf(err, "load %s", desc)
		}
	} else {
		desc = "default env conf"
	}

	env, err := env.New(json)
	if err != nil {
		return nil, errors.Maskf(err, "load %s", desc)
	}

	return env, nil
}

func parseMode() ([]string, error) {
	modes := make(map[string]bool)
	for _, m := range strings.Split(*mode, ",") {
		switch m {
		case "daemon", "initdb", "web":
			if modes[m] {
				return nil, errors.New("duplicate mode " + m)
			}
			modes[m] = true
		default:
			return nil, errors.New("unknown mode " + m)
		}
	}

	if modes["initdb"] && len(modes) > 1 {
		return nil, errors.New("initdb cannot be combined with other modes")
	}

	modesSlice := make([]string, len(modes))
	i := 0
	for m := range modes {
		modesSlice[i] = m
		i++
	}
	return modesSlice, nil
}

func initDB(e *env.Env) error {
	setupBytes, err := ioutil.ReadFile(SQL_SETUP_SCRIPT)
	if err != nil {
		return errors.Trace(err)
	}

	setupScript := string(setupBytes)

	err = e.DB.Setup(setupScript)

	return errors.Maskf(err, "unable to run setup script")
}

func waitForExitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGHUP, unix.SIGINT, unix.SIGTERM)
	<-c
}

func ifErrPrintAndExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
