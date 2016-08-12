package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type ResourceID int32
type ResourceType int16

type Resource struct {
	ResourceID   ResourceID
	ResourceType ResourceType
	Resource     string
}

func (db *DB) IsExistingResource(id ResourceID) (exists bool) {
	if id == 0 {
		return false
	}

	q := `SELECT EXISTS (SELECT * FROM pfx_resources WHERE resourceid = ?)`
	err := db.Get(&exists, cq(db, q), id)
	if err != nil {
		return false
	}
	return
}

func (tx *Tx) CreateResource(resource *Resource) (ResourceID, error) {
	q := `INSERT INTO pfx_resources (resourcetype, resource) VALUES (:resourcetype, :resource)`
	result, err := tx.NamedExec(cq(tx, q), resource)
	if err != nil {
		return 0, errors.Trace(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}
	return ResourceID(id), nil
}

func (tx *Tx) UpdateResource(resource *Resource) error {
	q := `UPDATE pfx_resources SET resourcetype=:resourcetype, resource=:resource
	     WHERE resourceid=:resourceid`
	_, err := tx.NamedExec(cq(tx, q), resource)
	return errors.Trace(err)
}

func (tx *Tx) LoadResource(id ResourceID) (*Resource, error) {
	var ds Resource
	q := `SELECT * FROM pfx_resources WHERE resourceid = ?`
	if err := tx.Get(&ds, cq(tx, q), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Trace(err)
	}
	return &ds, nil
}

func (db *DB) LoadResources() ([]*Resource, error) {
	var resources []*Resource
	q := `SELECT * FROM pfx_resources`
	if err := db.Select(&resources, cq(db, q)); err != nil {
		return nil, errors.Trace(err)
	}
	return resources, nil
}

func (db *DB) LoadResourcesOfTypes(ids []ResourceType) ([]*Resource, error) {
	var resources []*Resource
	q, args, err := sqlx.In(`SELECT * FROM pfx_resources WHERE resourcetype IN (?)`, ids)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if err := db.Select(&resources, cq(db, q), args...); err != nil {
		return nil, errors.Trace(err)
	}
	return resources, nil
}

func (tx *Tx) DeleteResource(id ResourceID) error {
	q := `DELETE FROM pfx_resources WHERE resourceid = ?`
	_, err := tx.Exec(cq(tx, q), id)
	return errors.Trace(err)
}
