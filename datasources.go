package revere

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere/datasources"
)

type DataSource struct {
	Id         uint                         `json:"id"`
	SourceType datasources.DataSourceTypeId `json:"sourceTypeId"`
	Source     string                       `json:"source"`
	Delete     bool                         `json:"delete,omitempty"`
}

const dataSourceFields = `id, sourceTypeId, source`

func (ds *DataSource) Validate() (errs []string) {
	dataSourceType, err := datasources.DataSourceTypeById(ds.SourceType)
	if err != nil {
		errs = append(errs, err.Error())
	}
	dataSource, err := dataSourceType.Load(ds.Source)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Invalid data source: %s", ds.Source))
	}
	errs = append(errs, dataSource.Validate()...)
	return
}

func LoadDataSourcesOfType(db *sql.DB, dataSourceTypeId datasources.DataSourceTypeId) ([]*DataSource, error) {
	return LoadDataSources(db, fmt.Sprintf("WHERE sourceTypeId = %d", dataSourceTypeId))
}

func LoadDataSourceById(db *sql.DB, id uint) (*DataSource, error) {
	results, err := LoadDataSources(db, fmt.Sprintf("WHERE id = %d", id))
	if len(results) == 0 {
		return nil, fmt.Errorf("Data source not found: %d", id)
	}
	return results[0], err
}

func LoadSourceContentOfType(db *sql.DB, allowedDataSource datasources.DataSourceTypeId) (sourceContent []string, err error) {
	sources, err := LoadDataSourcesOfType(db, allowedDataSource)
	if err != nil {
		return
	}

	for _, v := range sources {
		sourceContent = append(sourceContent, v.Source)
	}
	return
}

func LoadDataSources(db *sql.DB, condition string) (dataSources []*DataSource, err error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT %s FROM data_sources %s`, dataSourceFields, condition))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ds, err := loadDataSourceFromRow(rows)
		if err != nil {
			return nil, err
		}
		dataSources = append(dataSources, ds)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dataSources, nil
}

func (ds *DataSource) Save(db *sql.DB) (err error) {
	return Transact(db, func(tx *sql.Tx) error {
		if ds.isCreate() {
			_, err = ds.create(tx)
		} else if ds.Delete {
			err = ds.delete(tx)
		} else {
			err = ds.update(tx)
		}
		return err
	})
}

func (ds *DataSource) isCreate() bool {
	return ds.Id == 0
}

func (ds *DataSource) create(tx *sql.Tx) (uint, error) {
	stmt, err := tx.Prepare(`INSERT INTO data_sources (sourceTypeId, source) VALUES (?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ds.SourceType, ds.Source)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint(id), err
}

func (ds *DataSource) update(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE data_sources SET sourceTypeId = ?, source = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ds.SourceType, ds.Source, ds.Id)
	return err
}

func (ds *DataSource) delete(tx *sql.Tx) error {
	var stmt *sql.Stmt
	stmt, err := tx.Prepare(`
		DELETE FROM data_sources
		WHERE id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ds.Id)
	if err != nil {
		return err
	}
	return stmt.Close()
}

func loadDataSourceFromRow(rows *sql.Rows) (*DataSource, error) {
	var ds DataSource
	if err := rows.Scan(&ds.Id, &ds.SourceType, &ds.Source); err != nil {
		return nil, err
	}
	return &ds, nil
}
