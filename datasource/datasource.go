package datasource

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/web/tmpl"
)

// The VM struct is practically identical in purpose to its counterparts in the
// vm package, as it represents the intermediate structure between Revere's DB
// representation of the Data Source and its front end representation of the
// Data Source.
type VM struct {
	DataSource
	DataSourceParams string
	SourceType       db.SourceType
	SourceID         db.DatasourceID
	Delete           bool
}

// DataSourceType and DataSource define a common display abstraction for all
// Data Sources.
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
	MainScript    = "datasources.js"
)

var (
	defaultType = Graphite{}
	types       = make(map[db.SourceType]DataSourceType)
)

func Default() (DataSource, error) {
	ds, err := defaultType.blank()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return ds, nil
}

func AllScripts() []string {
	scripts := make([]string, 0)
	for _, dst := range types {
		for _, script := range dst.Scripts() {
			scripts = append(scripts, script)
		}
	}
	scripts = tmpl.AppendDir(DataSourceDir, scripts)
	scripts = append(scripts, MainScript)
	return scripts
}

func AllTypes() map[db.SourceType]string {
	typeIds := make(map[db.SourceType]string)
	for id, dst := range types {
		typeIds[id] = dst.Name()
	}
	return typeIds
}

func LoadFromParams(id db.SourceType, dsParams string) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return dsType.loadFromParams(dsParams)
}

func LoadFromDB(id db.SourceType, dsJson string) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return dsType.loadFromDB(dsJson)
}

func Blank(id db.SourceType) (DataSource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
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

func addType(dataSourceType DataSourceType) {
	if _, ok := types[dataSourceType.Id()]; !ok {
		types[dataSourceType.Id()] = dataSourceType
	} else {
		panic(fmt.Sprintf("A data source type with id %d already exists", dataSourceType.Id))
	}
}

func AllOfTypes(DB *db.DB, ids []db.SourceType) ([]*VM, error) {
	datasources, err := DB.LoadDatasourcesOfTypes(ids)
	if err != nil {
		return nil, errors.Trace(err)
	}

	dss, err := newVMs(datasources)

	return dss, errors.Trace(err)
}

func All(DB *db.DB) ([]*VM, error) {
	datasources, err := DB.LoadDatasources()
	if err != nil {
		return nil, errors.Trace(err)
	}

	dss, err := newVMs(datasources)

	return dss, errors.Trace(err)
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

func newVMs(datasources []*db.Datasource) ([]*VM, error) {
	dss := make([]*VM, len(datasources))
	var err error
	for i, datasource := range datasources {
		dss[i], err = newVM(datasource)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return dss, nil
}

func (vm *VM) Id() int64 {
	return int64(vm.SourceID)
}

func (*VM) ComponentName() string {
	return "Data Source"
}

func (vm *VM) IsCreate() bool {
	return vm.Id() == 0
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
	// Automatically short-circuit all deletions
	if vm.Delete {
		return nil
	}

	var err error
	vm.DataSource, err = LoadFromParams(vm.SourceType, vm.DataSourceParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load setting %s", vm.DataSourceParams))
		return errs
	}

	return vm.DataSource.Validate()
}
