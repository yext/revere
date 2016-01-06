package probes

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/yext/revere/util"
)

type ProbeTypeId int

type ProbeType interface {
	Id() ProbeTypeId
	Name() string
	Load(probe string) (Probe, error)
}

type Probe interface {
	Validate() []string
	Render() (template.HTML, error)
}

const probeTemplateDir = "web/views/probes/"

var (
	probeTypes map[ProbeTypeId]ProbeType = make(map[ProbeTypeId]ProbeType)

	// All probe templates
	tMap map[string]*template.Template = make(map[string]*template.Template)

	defaultProbeType     ProbeType = GraphiteThreshold{}
	defaultProbeTemplate template.HTML
)

func init() {
	// Fetch all probe templates
	funcMap := template.FuncMap{"strEq": util.StrEq}
	templateInfo, err := ioutil.ReadDir(probeTemplateDir)
	for _, t := range templateInfo {
		if t.IsDir() {
			continue
		}
		tMap[t.Name()], err = template.New(t.Name()).Funcs(funcMap).ParseFiles(probeTemplateDir + t.Name())
		if err != nil {
			panic(fmt.Sprintf("Got error initializing probe templates: %v", err))
		}
	}

	// Render the default probe template
	t, err := defaultProbeType.Load(`{}`)
	if err != nil {
		panic(fmt.Sprintf("Failed to load default probe template: %v", err))
	}

	template, err := t.Render()
	if err != nil {
		panic(fmt.Sprintf("Failed to render default probe template: %v", err))
	}
	defaultProbeTemplate = template
}

func GetProbeType(probeType ProbeTypeId) (ProbeType, error) {
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

func GetAllProbes() (pts []ProbeType) {
	for _, v := range probeTypes {
		pts = append(pts, v)
	}
	return pts
}

func LoadDefaultProbeTemplate() template.HTML {
	return defaultProbeTemplate
}
