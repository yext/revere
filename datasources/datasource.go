package datasources

import (
	"fmt"
)

type DataSourceTypeId int16

type DataSourceType interface {
	Id() DataSourceTypeId
	Name() string
	loadFromParams(ds string) (DataSource, error)
	loadFromDb(ds string) (DataSource, error)
	blank() (DataSource, error)
	Templates() string
	Scripts() []string
}

type DataSource interface {
	Serialize() (string, error)
	Type() DataSourceType
	Validate() []string
}

const (
	DataSourceDir = "datasources"
)

var (
	defaultType = Graphite{}
	types       = make(map[DataSourceTypeId]DataSourceType)
)

func Default() (DataSource, error) {
	ds, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func LoadFromParams(id DataSourceTypeId, dsParams string) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return dsType.loadFromParams(dsParams)
}

func LoadFromDb(id DataSourceTypeId, dsJson string) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return dsType.loadFromDb(dsJson)
}

func Blank(id DataSourceTypeId) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return dsType.blank()
}

func getType(id DataSourceTypeId) (DataSourceType, error) {
	dsType, ok := types[id]
	if !ok {
		return nil, fmt.Errorf("No data source type with id %d exists", id)
	}

	return dsType, nil
}

func addDataSourceType(dataSourceType DataSourceType) {
	if _, ok := types[dataSourceType.Id()]; !ok {
		types[dataSourceType.Id()] = dataSourceType
	} else {
		panic(fmt.Sprintf("A data source type with id %d already exists", dataSourceType.Id))
	}
}
