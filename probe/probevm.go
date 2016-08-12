package probe

import (
	"fmt"

	"github.com/yext/revere/db"
)

// ProbeVMType and ProbeVM define a common display abstraction for all probes.
type VMType interface {
	Id() db.ProbeType
	Name() string
	loadFromParams(probe string) (VM, error)
	loadFromDb(probe string, tx *db.Tx) (VM, error)
	blank() (VM, error)
	Templates() map[string]string
	Scripts() map[string][]string
	AcceptedResourceTypes() []db.ResourceType
}

type VM interface {
	VMType
	HasResource(db.ResourceID) bool
	SerializeForDB() (string, error)
	SerializeForFrontend() map[string]string
	Type() VMType
	Validate() []string
}

const (
	ProbesDir = "probes"
)

var (
	probeTypes  map[db.ProbeType]VMType = make(map[db.ProbeType]VMType)
	defaultType                         = GraphiteThresholdType{}
)

func Default() (VM, error) {
	probe, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return probe, nil
}

func LoadFromParams(id db.ProbeType, probeParams string) (VM, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromParams(probeParams)
}

func LoadFromDB(id db.ProbeType, probeJson string, tx *db.Tx) (VM, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromDb(probeJson, tx)
}

func Blank(id db.ProbeType) (VM, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}
	return probeType.blank()
}

func getType(id db.ProbeType) (VMType, error) {
	probeType, ok := probeTypes[id]
	if !ok {
		return nil, fmt.Errorf("No probe type with id %d exists", id)
	}

	return probeType, nil
}

func addType(probeType VMType) {
	if _, ok := probeTypes[probeType.Id()]; !ok {
		probeTypes[probeType.Id()] = probeType
	} else {
		panic(fmt.Sprintf("A probe type with id %d already exists", probeType.Id()))
	}
}

func AllTypes() (pts []VMType) {
	for _, v := range probeTypes {
		pts = append(pts, v)
	}
	return pts
}
