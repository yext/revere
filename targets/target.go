package targets

import "fmt"

type TargetTypeId int

type TargetType interface {
	Id() TargetTypeId
	Name() string
	Load(target string) (Target, error)
	Templates() map[string]string
	Scripts() map[string][]string
}

type Target interface {
	TargetType() TargetType
	Validate() (errs []string)
}

var (
	targetTypes map[TargetTypeId]TargetType = make(map[TargetTypeId]TargetType)
)

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
