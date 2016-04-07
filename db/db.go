// Package db encapsulates all interaction with Revere's backing database
// storage.
package db

import (
	"database/sql"
	"fmt"
	"strings"

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

// Prefix returns the prefix to add to table names in queries.
func (db *DB) Prefix() string {
	return db.prefix
}

func (db *DB) Unsafe() *DB {
	return &DB{DB: db.DB.Unsafe(), prefix: db.prefix}
}

func (db *DB) Beginx() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Tx{Tx: tx, prefix: db.prefix}, nil
}

type Tx struct {
	*sqlx.Tx
	prefix string
}

// Prefix returns the prefix to add to table names in queries.
func (tx *Tx) Prefix() string {
	return tx.prefix
}

func (tx *Tx) Unsafe() *Tx {
	return &Tx{Tx: tx.Tx.Unsafe(), prefix: tx.prefix}
}

// dbOrTx makes it easier to implement data loading methods that can be run
// either within or without a transaction.
type dbOrTx interface {
	// From sql.DB/Tx.

	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	// From sqlx.DB/Tx.

	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error

	// Custom to revere/db.

	Prefix() string
}

// cq customizes a query for issuing to the database. Query strings in this
// package are written with certain conventions that allow customizations to be
// automatically applied. For example, cq replaces all instances of pfx_ with
// the specific table prefix this instance of Revere has been configured with.
func cq(dt dbOrTx, query string) string {
	return strings.Replace(query, "pfx_", dt.Prefix(), -1)
}

func unsafe(dt dbOrTx) dbOrTx {
	switch dt := dt.(type) {
	case *DB:
		return dt.Unsafe()
	case *Tx:
		return dt.Unsafe()
	default:
		panic(fmt.Sprintf("revere/db.unsafe: unknown type: %T", dt))
	}
}
