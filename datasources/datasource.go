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
	LoadInfo(source string) (interface{}, error)
	DefaultInfo() interface{}
}

var (
	defaultDataSourceTypeId = GraphiteDataSource{}.Id()
	dataSourceTypes         = make(map[DataSourceTypeId]DataSourceType)
)

func DataSourceTypeById(dataSourceTypeId DataSourceTypeId) (DataSourceType, error) {
	if dst, ok := dataSourceTypes[dataSourceTypeId]; !ok {
		return dst, fmt.Errorf("Invalid data source type with id: %d", dataSourceTypeId)
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

func GetDataSourceTypes() (types []*DataSourceType) {
	for _, dst := range dataSourceTypes {
		types = append(types, &dst)
	}
	return
}

func DefaultDataSourceType() DataSourceType {
	return dataSourceTypes[defaultDataSourceTypeId]
}
