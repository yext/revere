// Package db encapsulates all interaction with Revere's backing database
// storage.
package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type DB struct {
	*sqlx.DB
	prefix string
}

// DBJSONModel provides the settings for a DB. It is used as the structure for
// configuring Revere's database connection in Revere's environment
// configuration JSON file.
type DBJSONModel struct {
	DSN         string
	TablePrefix string
}

// New validates conf and connects to the database specified in conf.
func New(conf DBJSONModel) (*DB, error) {
	// TODO(eefi): Validate conf.

	db, err := sqlx.Connect("mysql", conf.DSN)
	if err != nil {
		return nil, errors.Maskf(err, "connect")
	}

	// TODO(eefi): Validate DB has been initialized.

	return &DB{DB: db, prefix: conf.TablePrefix}, nil
}
