package datasources

import (
	"fmt"
)

type DataSourceTypeId uint

type DataSourceType interface {
	Id() DataSourceTypeId
	Name() string
	Template() string
	Scripts() []string
	Load(source string) (DataSource, error)
	LoadDefault() DataSource
}

type DataSource interface {
	Validate() []string
	DataSourceType() DataSourceType
}

var (
	defaultDataSourceTypeId = Graphite{}.Id()
	dataSourceTypes         = make(map[DataSourceTypeId]DataSourceType)
)

func DataSourceTypeById(dataSourceTypeId DataSourceTypeId) (DataSourceType, error) {
	if dst, ok := dataSourceTypes[dataSourceTypeId]; !ok {
		return nil, fmt.Errorf("Invalid data source type with id: %d", dataSourceTypeId)
	} else {
		return dst, nil
	}
}

func addDataSourceType(dataSourceType DataSourceType) {
	if _, ok := dataSourceTypes[dataSourceType.Id()]; !ok {
		dataSourceTypes[dataSourceType.Id()] = dataSourceType
	} else {
		panic(fmt.Sprintf("A data source type with id %d already exists", dataSourceType.Id))
	}
}

func GetDataSourceTypes() []*DataSourceType {
	types := make([]*DataSourceType, len(dataSourceTypes))
	i := 0
	for _, dst := range dataSourceTypes {
		types[i] = &dst
		i++
	}
	return types
}

func DefaultDataSourceType() DataSourceType {
	return dataSourceTypes[defaultDataSourceTypeId]
}
