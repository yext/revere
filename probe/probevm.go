package probe

import (
	"fmt"

	"github.com/yext/revere/db"
)

type ProbeVMType interface {
	Id() db.ProbeType
	Name() string
	loadFromParams(probe string) (ProbeVM, error)
	loadFromDb(probe string, tx *db.Tx) (ProbeVM, error)
	blank() (ProbeVM, error)
	Templates() map[string]string
	Scripts() map[string][]string
	AcceptedSourceTypes() []db.SourceType
}

type ProbeVM interface {
	ProbeVMType
	HasDatasource(db.DatasourceID) bool
	SerializeForDB() (string, error)
	SerializeForFrontend() map[string]string
	Type() ProbeVMType
	Validate() []string
}

const (
	ProbesDir = "probes"
)

var (
	probeTypes  map[db.ProbeType]ProbeVMType = make(map[db.ProbeType]ProbeVMType)
	defaultType                              = GraphiteThresholdType{}
)

func Default() (ProbeVM, error) {
	probe, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return probe, nil
}

func LoadFromParams(id db.ProbeType, probeParams string) (ProbeVM, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromParams(probeParams)
}

func LoadFromDB(id db.ProbeType, probeJson string, tx *db.Tx) (ProbeVM, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromDb(probeJson, tx)
}

func Blank(id db.ProbeType) (ProbeVM, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}
	return probeType.blank()
}

func getType(id db.ProbeType) (ProbeVMType, error) {
	probeType, ok := probeTypes[id]
	if !ok {
		return nil, fmt.Errorf("No probe type with id %d exists", id)
	}

	return probeType, nil
}

func addProbeVMType(probeType ProbeVMType) {
	if _, ok := probeTypes[probeType.Id()]; !ok {
		probeTypes[probeType.Id()] = probeType
	} else {
		panic(fmt.Sprintf("A probe type with id %d already exists", probeType.Id()))
	}
}

func AllTypes() (pts []ProbeVMType) {
	for _, v := range probeTypes {
		pts = append(pts, v)
	}
	return pts
}
