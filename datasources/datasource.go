package datasources

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
)

type VM struct {
	DataSource
	DataSourceParams string
	SourceType       db.SourceType
	SourceID         db.DatasourceID
	Delete           bool
}

type DataSourceType interface {
	Id() db.SourceType
	Name() string
	loadFromParams(ds string) (DataSource, error)
	loadFromDB(ds string) (DataSource, error)
	blank() (DataSource, error)
	Templates() string
	Scripts() []string
}

type DataSource interface {
	DataSourceType
	Serialize() (string, error)
	Type() DataSourceType
	Validate() []string
}

const (
	DataSourceDir = "datasources"
)

var (
	defaultType = Graphite{}
	types       = make(map[db.SourceType]DataSourceType)
)

func Default() (DataSource, error) {
	ds, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func LoadFromParams(id db.SourceType, dsParams string) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return dsType.loadFromParams(dsParams)
}

func LoadFromDB(id db.SourceType, dsJson string) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return dsType.loadFromDB(dsJson)
}

func Blank(id db.SourceType) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return dsType.blank()
}

func getType(id db.SourceType) (DataSourceType, error) {
	dsType, ok := types[id]
	if !ok {
		return nil, fmt.Errorf("No data source type with id %d exists", id)
	}

	return dsType, nil
}

//TODO (fchen): eventually change to addType for datasources, settings, probes, targets; conform naming and syntax
func addDataSourceType(dataSourceType DataSourceType) {
	if _, ok := types[dataSourceType.Id()]; !ok {
		types[dataSourceType.Id()] = dataSourceType
	} else {
		panic(fmt.Sprintf("A data source type with id %d already exists", dataSourceType.Id))
	}
}

func AllTypes() (dsts []DataSourceType) {
	for _, t := range types {
		dsts = append(dsts, t)
	}
	return dsts
}

func All(DB *db.DB) ([]*VM, error) {
	datasources, err := DB.LoadDatasources()
	if err != nil {
		return nil, err
	}

	dss := make([]*VM, len(datasources))
	for i, datasource := range datasources {
		dss[i], err = newVM(datasource)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return dss, nil
}

func newVM(ds *db.Datasource) (*VM, error) {
	datasource, err := LoadFromDB(ds.SourceType, ds.Source)
	if err != nil {
		return &VM{}, errors.Trace(err)
	}

	return &VM{
		DataSource: datasource,
		SourceID:   ds.SourceID,
		SourceType: ds.SourceType,
	}, nil
}

func (vm *VM) IsCreate() bool {
	return vm.SourceID == 0
}

func (vm *VM) IsDelete() bool {
	return vm.Delete
}

func (vm *VM) Save(tx *db.Tx) error {
	var err error
	vm.DataSource, err = LoadFromParams(vm.SourceType, vm.DataSourceParams)
	if err != nil {
		return errors.Trace(err)
	}

	dsJSON, err := vm.DataSource.Serialize()
	if err != nil {
		return errors.Trace(err)
	}

	datasource := &db.Datasource{
		SourceID:   vm.SourceID,
		SourceType: vm.SourceType,
		Source:     dsJSON,
	}

	if vm.IsCreate() {
		var id db.DatasourceID
		id, err = tx.CreateDatasource(datasource)
		datasource.SourceID = id
	} else if vm.IsDelete() {
		err = tx.DeleteDatasource(vm.SourceID)
	} else {
		err = tx.UpdateDatasource(datasource)
	}

	return errors.Trace(err)
}

func (vm *VM) Validate() (errs []string) {
	var err error
	vm.DataSource, err = LoadFromParams(vm.SourceType, vm.DataSourceParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load setting %s", vm.DataSourceParams))
		return errs
	}

	return vm.DataSource.Validate()
}
