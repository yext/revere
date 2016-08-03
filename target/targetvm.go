package target

import (
	"fmt"

	"github.com/yext/revere/db"
)

type TargetTypeVM interface {
	Id() db.TargetType
	Name() string
	loadFromParams(target string) (TargetVM, error)
	loadFromDb(target string) (TargetVM, error)
	blank() TargetVM
	Templates() map[string]string
	Scripts() map[string][]string
}

type TargetVM interface {
	TargetTypeVM
	Serialize() (string, error)
	Type() TargetTypeVM
	Validate() []string
}

const (
	TargetsDir = "targets"
)

var (
	targetTypes map[db.TargetType]TargetTypeVM = make(map[db.TargetType]TargetTypeVM)
	defaultType                                = EmailType{}
)

func Default() TargetVM {
	return defaultType.blank()
}

func LoadFromParams(id db.TargetType, targetParams string) (TargetVM, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromParams(targetParams)
}

func LoadFromDb(id db.TargetType, targetJson string) (TargetVM, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromDb(targetJson)
}

func Blank(id db.TargetType) (TargetVM, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.blank(), nil
}

func getType(id db.TargetType) (TargetTypeVM, error) {
	targetType, ok := targetTypes[id]
	if !ok {
		return nil, fmt.Errorf("No target type with id %d exists", id)
	}

	return targetType, nil
}

func addTargetVMType(targetType TargetTypeVM) {
	if _, ok := targetTypes[targetType.Id()]; !ok {
		targetTypes[targetType.Id()] = targetType
	} else {
		panic(fmt.Sprintf("A target type with id %d already exists", targetType.Id()))
	}
}

func AllTargets() (tts []TargetTypeVM) {
	for _, v := range targetTypes {
		tts = append(tts, v)
	}
	return tts
}
