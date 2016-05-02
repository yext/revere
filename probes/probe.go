package probes

import (
	"fmt"

	"github.com/yext/revere/datasources"
	"github.com/yext/revere/probes"
)

type ProbeTypeId int16

type ProbeType interface {
	Id() ProbeTypeId
	Name() string
	loadFromParams(probe string) (Probe, error)
	loadFromDb(probe string) (Probe, error)
	blank() (Probe, error)
	Templates() map[string]string
	Scripts() map[string][]string
	AcceptedDataSourceTypeIds() []datasources.DataSourceTypeId
}

type Probe interface {
	Serialize() (string, error)
	ProbeType() ProbeType
	Validate() []string
}

const (
	ProbesDir = "probes"
)

var (
	types       map[ProbeTypeId]ProbeType = make(map[ProbeTypeId]ProbeType)
	defaultType                           = probes.GraphiteThreshold{}
)

func Default() (Probe, error) {
	probe, err := Blank(defaultType)
	if err != nil {
		return nil, err
	}

	return probe, nil
}

func LoadFromParams(id ProbeTypeId, probeParams string) (Probe, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromParams(probeParams)
}

func LoadFromDb(id ProbeTypeId, probeJson string) (Probe, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.loadFromDb(probeJson)
}

func Blank(id ProbeTypeId) (Probe, error) {
	probeType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return probeType.blank()
}

func getType(id ProbeTypeId) (ProbeType, error) {
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
