package vm

import (
	"database/sql"
	"fmt"

	"github.com/yext/revere"
	"github.com/yext/revere/datasources"
	"github.com/yext/revere/probes"
)

type Probe struct {
	probes.Probe
	DataSourceInfo map[string][]interface{}
}

var (
	defaultProbeTypeId = probes.GraphiteThreshold{}.Id()
)

const (
	ProbesDir = "probes"
)

func NewProbe(db *sql.DB, p probes.Probe) (*Probe, error) {
	viewmodel := new(Probe)
	viewmodel.Probe = p
	viewmodel.DataSourceInfo = make(map[string][]interface{})
	for _, id := range p.ProbeType().AcceptedDataSourceTypeIds() {
		dst, _ := datasources.DataSourceTypeById(id)
		dataSourceInfosOfType := make([]interface{}, 0)
		sourceContent, err := revere.LoadSourceContentOfType(db, id)
		if err != nil {
			return nil, err
		}
		for _, s := range sourceContent {
			info, err := dst.LoadInfo(s)
			if err != nil {
				return nil, err
			}
			dataSourceInfosOfType = append(dataSourceInfosOfType, info)
		}
		viewmodel.DataSourceInfo[dst.Name()] = dataSourceInfosOfType
	}

	return viewmodel, nil
}

func DefaultProbe(db *sql.DB) *Probe {
	probe, err := BlankProbe(db, int(defaultProbeTypeId))
	if err != nil {
		panic(err)
	}

	return probe
}

func BlankProbe(db *sql.DB, pt int) (*Probe, error) {
	probeType, err := probes.ProbeTypeById(probes.ProbeTypeId(pt))
	if err != nil {
		return nil, fmt.Errorf("Probe type not found: %d", pt)
	}

	probe, err := probeType.Load(`{}`)
	if err != nil {
		return nil, fmt.Errorf("Unable to load %s", probeType.Name())
	}

	return NewProbe(db, probe)
}
