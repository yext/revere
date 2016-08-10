package target

import (
	"fmt"

	"github.com/yext/revere/db"
)

// TargetVMType and TargetVM define a common display abstraction for all
// targets.
type VMType interface {
	Id() db.TargetType
	Name() string
	loadFromParams(target string) (VM, error)
	loadFromDb(target string) (VM, error)
	blank() VM
	Templates() map[string]string
	Scripts() map[string][]string
}

type VM interface {
	VMType
	Serialize() (string, error)
	Type() VMType
	Validate() []string
}

const (
	TargetsDir = "targets"
)

var (
	targetTypes map[db.TargetType]VMType = make(map[db.TargetType]VMType)
	defaultType                          = EmailType{}
)

func Default() VM {
	return defaultType.blank()
}

func LoadFromParams(id db.TargetType, targetParams string) (VM, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromParams(targetParams)
}

func LoadFromDb(id db.TargetType, targetJson string) (VM, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.loadFromDb(targetJson)
}

func Blank(id db.TargetType) (VM, error) {
	targetType, err := getType(id)
	if err != nil {
		return nil, err
	}

	return targetType.blank(), nil
}

func getType(id db.TargetType) (VMType, error) {
	targetType, ok := targetTypes[id]
	if !ok {
		return nil, fmt.Errorf("No target type with id %d exists", id)
	}

	return targetType, nil
}

func addType(targetType VMType) {
	if _, ok := targetTypes[targetType.Id()]; !ok {
		targetTypes[targetType.Id()] = targetType
	} else {
		panic(fmt.Sprintf("A target type with id %d already exists", targetType.Id()))
	}
}

func AllTargets() (tts []VMType) {
	for _, v := range targetTypes {
		tts = append(tts, v)
	}
	return tts
}
