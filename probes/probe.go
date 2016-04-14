package probes

import (
	"fmt"

	"github.com/yext/revere/datasources"
)

type ProbeTypeId int16

type ProbeType interface {
	Id() ProbeTypeId
	Name() string
	Load(probe string) (Probe, error)
	Templates() map[string]string
	Scripts() map[string][]string
	AcceptedDataSourceTypeIds() []datasources.DataSourceTypeId
}

type Probe interface {
	ProbeType() ProbeType
	Validate() []string
}

var (
	probeTypes map[ProbeTypeId]ProbeType = make(map[ProbeTypeId]ProbeType)
)

func init() {
}

func ProbeTypeById(probeType ProbeTypeId) (ProbeType, error) {
	if pt, ok := probeTypes[probeType]; !ok {
		return pt, fmt.Errorf("Invalid probe type with id: %d", probeType)
	} else {
		return pt, nil
	}
}

func addProbeType(probeType ProbeType) {
	if _, ok := probeTypes[probeType.Id()]; !ok {
		probeTypes[probeType.Id()] = probeType
	} else {
		panic(fmt.Sprintf("A probe type with id %d already exists", probeType.Id()))
	}
}

func AllProbes() (pts []ProbeType) {
	for _, v := range probeTypes {
		pts = append(pts, v)
	}
	return pts
}
