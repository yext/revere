package vm

import (
	"github.com/yext/revere"
	"github.com/yext/revere/datasources"
)

type DataSource struct {
	*revere.DataSource
	Attributes     datasources.DataSource
	DataSourceType datasources.DataSourceType
}

func NewDataSourceViewModel(ds *revere.DataSource) (*DataSource, error) {
	viewmodel := new(DataSource)

	viewmodel.DataSource = ds

	dataSourceType, err := datasources.DataSourceTypeById(ds.SourceType)
	if err != nil {
		return nil, err
	}
	viewmodel.Attributes, err = dataSourceType.Load(ds.Source)
	if err != nil {
		return nil, err
	}

	viewmodel.DataSourceType = dataSourceType
	return viewmodel, nil
}

func BlankDataSourceViewModelWithType(typeId datasources.DataSourceTypeId) (viewmodel *DataSource, err error) {
	viewmodel = new(DataSource)
	viewmodel.DataSource = new(revere.DataSource)
	viewmodel.DataSourceType, err = datasources.DataSourceTypeById(typeId)
	if err != nil {
		viewmodel = nil
		return
	}
	viewmodel.Attributes = viewmodel.DataSourceType.LoadDefault()

	return
}
