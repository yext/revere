package vm

import (
	"fmt"

	"github.com/yext/revere/targets"
)

type Target struct {
	targets.Target
}

var (
	defaultTargetTypeId = targets.Email{}.Id()
)

const (
	TargetsDir = "targets"
)

func NewTarget(t targets.Target) *Target {
	viewmodel := new(Target)
	viewmodel.Target = t

	return viewmodel
}

func DefaultTarget() *Target {
	target, err := BlankTarget(int(defaultTargetTypeId))
	if err != nil {
		panic(err)
	}

	return target
}

func BlankTarget(tt int) (*Target, error) {
	targetType, err := targets.TargetTypeById(targets.TargetTypeId(tt))
	if err != nil {
		return nil, fmt.Errorf("Target type not found: %d", tt)
	}

	target, err := targetType.Load(`{}`)
	if err != nil {
		return nil, fmt.Errorf("Unable to load %s", targetType.Name())
	}

	return NewTarget(target), nil
}
