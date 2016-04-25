package vm

import (
	"database/sql"
	"github.com/yext/revere"
	"github.com/yext/revere/datasources"
)

type DataSourceType struct {
	Type        datasources.DataSourceType
	DataSources []*DataSource
}

const (
	DataSourceDir = "datasources"
)

func NewDataSourceTypeViewModel(db *sql.DB, dst *datasources.DataSourceType) (*DataSourceType, error) {
	dstvm := new(DataSourceType)
	dstvm.Type = *dst

	dataSources, err := revere.LoadDataSourcesOfType(db, dstvm.Type.Id())
	if err != nil {
		return nil, err
	}
	arr := make([]*DataSource, 0)
	for _, ds := range dataSources {
		new, err := NewDataSource(ds)
		if err != nil {
			return nil, err
		}
		arr = append(arr, new)
	}
	if len(arr) == 0 {
		new, err := BlankDataSource(dstvm.Type.Id())
		if err != nil {
			return nil, err
		}
		arr = append(arr, new)
	}

	dstvm.DataSources = arr
	return dstvm, nil
}
