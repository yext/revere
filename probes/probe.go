package probes

import (
	"fmt"
	"html/template"

	"github.com/yext/revere/util"
)

type ProbeTypeId int

type ProbeType interface {
	Id() ProbeTypeId
	Name() string
	Load(probe string) (Probe, error)
}

type Probe interface {
	ProbeType() ProbeType
	Validate() []string
	Render() (template.HTML, error)
}

const probeTemplateDir = "web/views/probes/"

var (
	probeTypes map[ProbeTypeId]ProbeType = make(map[ProbeTypeId]ProbeType)

	// All probe templates
	probeTemplates map[string]*template.Template

	defaultProbeType ProbeType = GraphiteThreshold{}
)

func init() {
	// Fetch all probe templates
	probeTemplates = util.InitTemplates(probeTemplateDir, template.FuncMap{"strEq": util.StrEq})
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

func DefaultProbeTemplate() (template.HTML, error) {
	// Render the default probe template
	t, err := defaultProbeType.Load(`{}`)
	if err != nil {
		return "", err
	}

	template, err := t.Render()
	if err != nil {
		return "", err
	}
	return template, nil
}
