// Package env handles Revere's static environment. This environment is
// configured via a JSON file provided on the command line when running Revere,
// and it does not change at runtime.
package env

import (
	"encoding/json"

	"github.com/juju/errors"

	"github.com/yext/revere/db"
)

// Env provides runtime access to Revere's static environment.
type Env struct {
	DB   *db.DB
	Port uint16
	Host string
}

// New initializes an Env based on the configuration found in conf, which
// contains a serialized JSON object.
func New(conf []byte) (*Env, error) {
	var model EnvJSONModel
	err := json.Unmarshal(conf, &model)
	if err != nil {
		switch err := err.(type) {
		case *json.SyntaxError:
			return nil, errors.Maskf(
				err, "parse at byte %d", err.Offset)
		case *json.UnmarshalTypeError:
			return nil, errors.Maskf(
				err, "parse at byte %d", err.Offset)
		default:
			return nil, errors.Maskf(err, "parse")
		}
	}

	var e Env
	e.DB, err = db.New(model.DB)
	if err != nil {
		return nil, errors.Maskf(err, "load DB")
	}
	e.Port = model.Port
	e.Host = model.Host

	return &e, nil
}

// EnvJSONModel is the structure for Revere's environment configuration JSON
// file.
type EnvJSONModel struct {
	DB   db.DBJSONModel
	Port uint16
	Host string
}
