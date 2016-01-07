package targets

import (
	"fmt"
	"html/template"

	"github.com/yext/revere/util"
)

type Target interface {
	Validate() (errs []string)
	Render() (template.HTML, error)
}

type TargetTypeId int

type TargetType interface {
	Id() TargetTypeId
	Name() string
	Load(target string) (Target, error)
}

var (
	targetTemplateDir = "web/views/targets/"

	targetTypes map[TargetTypeId]TargetType = make(map[TargetTypeId]TargetType)

	targetTemplates map[string]*template.Template

	defaultTargetType     TargetType = Email{}
	defaultTargetTemplate template.HTML
)

func init() {
	targetTemplates = util.InitTemplates(targetTemplateDir, template.FuncMap{"strEq": util.StrEq})

	t, err := defaultTargetType.Load(`{}`)
	if err != nil {
		panic(fmt.Sprintf("Failed to load default target type: %v", err))
	}

	template, err := t.Render()
	if err != nil {
		panic(fmt.Sprintf("Failed to render default target type: %v", err))
	}
	defaultTargetTemplate = template
}

func TargetTypeById(targetType TargetTypeId) (TargetType, error) {
	if tt, ok := targetTypes[targetType]; !ok {
		return tt, fmt.Errorf("Invalid target type %d", targetType)
	} else {
		return tt, nil
	}
}

func addTargetType(targetType TargetType) {
	if _, ok := targetTypes[targetType.Id()]; !ok {
		targetTypes[targetType.Id()] = targetType
	} else {
		panic(fmt.Sprintf("A target type with id %d already exists", targetType.Id()))
	}
}

func AllTargets() (tts []TargetType) {
	for _, v := range targetTypes {
		tts = append(tts, v)
	}
	return tts
}

func DefaultTargetTemplate() template.HTML {
	return defaultTargetTemplate
}
