package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type DatasourceID int32
type SourceType int16

type Datasource struct {
	SourceID   DatasourceID
	SourceType SourceType
	Source     string
}

func (db *DB) IsExistingDatasource(id DatasourceID) (exists bool) {
	if id == 0 {
		return false
	}

	q := `SELECT EXISTS (SELECT * FROM pfx_data_sources WHERE sourceid = ?)`
	err := db.Get(&exists, cq(db, q), id)
	if err != nil {
		return false
	}
	return
}

func (tx *Tx) CreateDatasource(datasource *Datasource) (DatasourceID, error) {
	q := `INSERT INTO pfx_data_sources (sourcetype, source) VALUES (:sourcetype, :source)`
	result, err := tx.NamedExec(cq(tx, q), datasource)
	if err != nil {
		return 0, errors.Trace(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}
	return DatasourceID(id), nil
}

func (tx *Tx) UpdateDatasource(datasource *Datasource) error {
	q := `UPDATE pfx_data_sources SET sourcetype=:sourcetype, source=:source
	     WHERE sourceid=:sourceid`
	_, err := tx.NamedExec(cq(tx, q), datasource)
	return errors.Trace(err)
}

func (db *DB) LoadDatasource(id DatasourceID) (*Datasource, error) {
	var ds Datasource
	q := `SELECT * FROM pfx_data_sources WHERE sourceid = ?`
	if err := db.Get(&ds, cq(db, q), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &ds, nil
}

func (db *DB) LoadDatasources() ([]*Datasource, error) {
	var datasources []*Datasource
	q := `SELECT * FROM pfx_data_sources`
	if err := db.Select(&datasources, cq(db, q)); err != nil {
		return nil, errors.Trace(err)
	}
	return datasources, nil
}

func (db *DB) LoadDatasourcesOfTypes(ids []SourceType) ([]*Datasource, error) {
	var datasources []*Datasource
	q, args, err := sqlx.In(`SELECT * FROM pfx_data_sources WHERE sourcetype IN (?)`, ids)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if err := db.Select(&datasources, cq(db, q), args...); err != nil {
		return nil, errors.Trace(err)
	}
	return datasources, nil
}

func (tx *Tx) DeleteDatasource(id DatasourceID) error {
	q := `DELETE FROM pfx_data_sources WHERE sourceid = ?`
	_, err := tx.Exec(cq(tx, q), id)
	return errors.Trace(err)
}
