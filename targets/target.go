package targets

import (
	"fmt"

	"github.com/yext/revere/db"
)

type TargetType interface {
	Id() db.TargetType
	Name() string
	loadFromParams(target string) (Target, error)
	loadFromDb(target string) (Target, error)
	blank() (Target, error)
	Templates() map[string]string
	Scripts() map[string][]string
}

type Target interface {
	TargetType
	Serialize() (string, error)
	Type() TargetType
	Validate() []string
}

const (
	TargetsDir = "targets"
)

var (
	types       map[db.TargetType]TargetType = make(map[db.TargetType]TargetType)
	defaultType                              = Email{}
)

func Default() (Target, error) {
	target, err := defaultType.blank()
	if err != nil {
		return nil, err
	}

	return target, nil
}

func LoadFromParams(id db.TargetType, targetParams string) (Target, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromParams(targetParams)
}

func LoadFromDb(id db.TargetType, targetJson string) (Target, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromDb(targetJson)
}

func Blank(id db.TargetType) (Target, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.blank()
}

func getType(id db.TargetType) (TargetType, error) {
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
