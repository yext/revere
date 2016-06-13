package probes

import (
	"fmt"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/db"
)

type ProbeType interface {
	Id() db.ProbeType
	Name() string
	loadFromParams(probe string) (Probe, error)
	loadFromDb(probe string) (Probe, error)
	blank() (Probe, error)
	Templates() map[string]string
	Scripts() map[string][]string
	AcceptedDataSourceTypeIds() []datasources.DataSourceTypeId
}

type Probe interface {
	ProbeType
	Serialize() (string, error)
	Type() ProbeType
	Validate() []string
}

const (
	ProbesDir = "probes"
)

var (
	types       map[db.ProbeType]ProbeType = make(map[db.ProbeType]ProbeType)
	defaultType                            = GraphiteThreshold{}
)

func Default() (Probe, error) {
	probe, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return probe, nil
}

func LoadFromParams(id db.ProbeType, probeParams string) (Probe, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromParams(probeParams)
}

func LoadFromDB(id db.ProbeType, probeJson string) (Probe, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromDb(probeJson)
}

func Blank(id db.ProbeType) (Probe, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}
	return probeType.blank()
}

func getType(id db.ProbeType) (ProbeType, error) {
	probeType, ok := types[id]
	if !ok {
		return nil, fmt.Errorf("No probe type with id %d exists", id)
	}

	return probeType, nil
}

func addProbeType(probeType ProbeType) {
	if _, ok := types[probeType.Id()]; !ok {
		types[probeType.Id()] = probeType
	} else {
		panic(fmt.Sprintf("A probe type with id %d already exists", probeType.Id()))
	}
}

func AllTypes() (pts []ProbeType) {
	for _, v := range types {
		pts = append(pts, v)
	}
	return pts
}
