package targets

import "fmt"

type TargetTypeId int16

type TargetType interface {
	Id() TargetTypeId
	Name() string
	loadFromParams(target string) (Target, error)
	loadFromDb(target string) (Target, error)
	blank() (Target, error)
	Templates() map[string]string
	Scripts() map[string][]string
}

type Target interface {
	Serialize() (string, error)
	Type() TargetType
	Validate() []string
}

const (
	TargetsDir = "targets"
)

var (
	types       map[TargetTypeId]TargetType = make(map[TargetTypeId]TargetType)
	defaultType                             = Email{}
)

func Default() (Target, error) {
	target, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return target, nil
}

func LoadFromParams(id TargetTypeId, targetParams string) (Target, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromParams(targetParams)
}

func LoadFromDb(id TargetTypeId, targetJson string) (Target, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromDb(targetJson)
}

func Blank(id TargetTypeId) (Target, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.blank()
}

func getType(id TargetTypeId) (TargetType, error) {
	targetType, ok := types[id]
	if !ok {
		return nil, fmt.Errorf("No target type with id %d exists", id)
	}

	return targetType, nil
}

func addTargetType(targetType TargetType) {
	if _, ok := types[targetType.Id()]; !ok {
		types[targetType.Id()] = targetType
	} else {
		panic(fmt.Sprintf("A target type with id %d already exists", targetType.Id()))
	}
}

func AllTargets() (tts []TargetType) {
	for _, v := range types {
		tts = append(tts, v)
	}
	return tts
}
