package revere

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/juju/errgo"
)

type Env struct {
	db *sql.DB

	// web
	port    int
	urlBase string

	// Graphite data sources
	graphiteUrls []string
}

func (e *Env) Db() *sql.DB {
	return e.db
}

func (e *Env) Port() int {
	return e.port
}

func (e *Env) GraphiteUrls() []string {
	return e.graphiteUrls
}

type jsonConf struct {
	Db       *jsonDbConf
	Web      *jsonWebConf
	Graphite *jsonGraphiteConf
}

type jsonDbConf struct {
	Host     string
	Port     int
	Username string
	Password string
	Dsn      string
	// TODO(dp): use this
	TablePrefix string
}

// TODO(dp): use this
type jsonWebConf struct {
	Port    int
	UrlBase string
}

// TODO(psingh): Remove once we put data sources in the db
type jsonGraphiteConf struct {
	Urls []string
}

var DEFAULT_CONF_FILE = "conf/defaults.conf"

func BuildEnvFromFile(fn string) (*Env, error) {
	confBytes, err := ioutil.ReadFile(DEFAULT_CONF_FILE)
	if err != nil {
		fmt.Printf("Unable to find default conf file: %s. This file shouldn't be moved or edited.\n", DEFAULT_CONF_FILE)
		return nil, errgo.Mask(err)
	}

	var conf jsonConf
	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		fmt.Printf("Unable to parse default conf file: %s. This file shouldn't be edited.\n", DEFAULT_CONF_FILE)
		return nil, errgo.Mask(err)
	}

	confBytes, err = ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("Unable to find your conf file: %s. Double check the path.\n", fn)
		return nil, errgo.Mask(err)
	}

	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		fmt.Printf("Unable to parse your conf file: %s. Validate the json configuration.\n", fn)
		return nil, errgo.Mask(err)
	}

	var env Env
	env.db, err = getDb(conf.Db)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	env.port = conf.Web.Port
	env.urlBase = conf.Web.UrlBase
	env.graphiteUrls = conf.Graphite.Urls

	return &env, nil
}

// TODO(dp): test this
// getDb returns a sql db connection if the configuration is valid
func getDb(conf *jsonDbConf) (*sql.DB, error) {
	dsn, err := getValidDsn(conf.Dsn)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	if dsn == "" {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/revere?loc=UTC&parseTime=true",
			conf.Username, conf.Password, conf.Host, conf.Port)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Unable to connect to mysql with dsn: %s. Check your conf file.\n", dsn)
		return nil, errgo.Mask(err)
	}

	return db, nil
}

// getDsn validates a dsn string and returns it or an error if it's invalid
func getValidDsn(d string) (string, error) {
	if d == "" {
		return "", nil
	}
	params, err := url.ParseQuery(strings.SplitN(d, "?", 2)[1])
	if err != nil {
		fmt.Printf("For dsn %s, query params appeared invalid. Check anything after the ?.\n", d)
		return "", errgo.Mask(err)
	}

	if len(params["parseTime"]) > 1 || params.Get("parseTime") != "true" {
		fmt.Printf("For dsn %s, required parseTime to be true, but was %s\n", d, params.Get("parseTime"))
		return "", errgo.Newf("Required parseTime to be true, but was %s", params.Get("parseTime"))
	}

	if len(params["loc"]) > 1 || (params.Get("loc") != "UTC" && params.Get("loc") != "") {
		fmt.Printf("For dsn %s, required loc to be UTC or unset, but was %s\n", d, params.Get("loc"))
		return "", errgo.Mask(errors.New(fmt.Sprintf("Required loc to be UTC or unset, but was %s", params.Get("loc"))))
	}

	return d, nil
}
