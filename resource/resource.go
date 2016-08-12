package resource

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/yext/revere/db"
	"github.com/yext/revere/web/tmpl"
)

// The VM struct is practically identical in purpose to its counterparts in the
// vm package, as it represents the intermediate structure between Revere's DB
// representation of the Resource and its front end representation of the
// Resource.
type VM struct {
	Resource
	ResourceParams string
	ResourceType   db.ResourceType
	ResourceID     db.ResourceID
	Delete         bool
}

// ResourceType and Resource define a common display abstraction for all
// Resources.
type ResourceType interface {
	Id() db.ResourceType
	Name() string
	loadFromParams(ds string) (Resource, error)
	loadFromDB(ds string) (Resource, error)
	blank() (Resource, error)
	Templates() string
	Scripts() []string
}

type Resource interface {
	ResourceType
	Serialize() (string, error)
	Type() ResourceType
	Validate() []string
}

const (
	ResourceDir = "resources"
	MainScript  = "resources.js"
)

var (
	defaultType = Graphite{}
	types       = make(map[db.ResourceType]ResourceType)
)

func Default() (Resource, error) {
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
	scripts = tmpl.AppendDir(ResourceDir, scripts)
	scripts = append(scripts, MainScript)
	return scripts
}

func AllTypes() map[db.ResourceType]string {
	typeIds := make(map[db.ResourceType]string)
	for id, dst := range types {
		typeIds[id] = dst.Name()
	}
	return typeIds
}

func LoadFromParams(id db.ResourceType, dsParams string) (Resource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return dsType.loadFromParams(dsParams)
}

func LoadFromDB(id db.ResourceType, dsJson string) (Resource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return dsType.loadFromDB(dsJson)
}

func Blank(id db.ResourceType) (Resource, error) {
	dsType, err := getType(id)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return dsType.blank()
}

func getType(id db.ResourceType) (ResourceType, error) {
	dsType, ok := types[id]
	if !ok {
		return nil, fmt.Errorf("No resource type with id %d exists", id)
	}

	return dsType, nil
}

func addType(resourceType ResourceType) {
	if _, ok := types[resourceType.Id()]; !ok {
		types[resourceType.Id()] = resourceType
	} else {
		panic(fmt.Sprintf("A resource type with id %d already exists", resourceType.Id))
	}
}

func AllOfTypes(DB *db.DB, ids []db.ResourceType) ([]*VM, error) {
	resources, err := DB.LoadResourcesOfTypes(ids)
	if err != nil {
		return nil, errors.Trace(err)
	}

	dss, err := newVMs(resources)

	return dss, errors.Trace(err)
}

func All(DB *db.DB) ([]*VM, error) {
	resources, err := DB.LoadResources()
	if err != nil {
		return nil, errors.Trace(err)
	}

	dss, err := newVMs(resources)

	return dss, errors.Trace(err)
}

func newVM(ds *db.Resource) (*VM, error) {
	resource, err := LoadFromDB(ds.ResourceType, ds.Resource)
	if err != nil {
		return &VM{}, errors.Trace(err)
	}

	return &VM{
		Resource:     resource,
		ResourceID:   ds.ResourceID,
		ResourceType: ds.ResourceType,
	}, nil
}

func newVMs(resources []*db.Resource) ([]*VM, error) {
	dss := make([]*VM, len(resources))
	var err error
	for i, resource := range resources {
		dss[i], err = newVM(resource)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return dss, nil
}

func (vm *VM) Id() int64 {
	return int64(vm.ResourceID)
}

func (*VM) ComponentName() string {
	return "Resource"
}

func (vm *VM) IsCreate() bool {
	return vm.Id() == 0
}

func (vm *VM) IsDelete() bool {
	return vm.Delete
}

func (vm *VM) Save(tx *db.Tx) error {
	var err error
	vm.Resource, err = LoadFromParams(vm.ResourceType, vm.ResourceParams)
	if err != nil {
		return errors.Trace(err)
	}

	dsJSON, err := vm.Resource.Serialize()
	if err != nil {
		return errors.Trace(err)
	}

	resource := &db.Resource{
		ResourceID:   vm.ResourceID,
		ResourceType: vm.ResourceType,
		Resource:     dsJSON,
	}

	if vm.IsCreate() {
		var id db.ResourceID
		id, err = tx.CreateResource(resource)
		resource.ResourceID = id
	} else if vm.IsDelete() {
		err = tx.DeleteResource(vm.ResourceID)
	} else {
		err = tx.UpdateResource(resource)
	}

	return errors.Trace(err)
}

func (vm *VM) Validate() (errs []string) {
	var err error
	vm.Resource, err = LoadFromParams(vm.ResourceType, vm.ResourceParams)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Unable to load setting %s", vm.ResourceParams))
		return errs
	}

	return vm.Resource.Validate()
}
