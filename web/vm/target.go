package vm

import (
	"fmt"

	"github.com/yext/revere/targets"
)

type Target struct {
	targets.Target
	templates map[string]string
	scripts   map[string][]string
}

var (
	defaultTargetTypeId = targets.Email{}.Id()
)

const (
	targetsDir = "targets"
)

func NewTarget(t targets.Target) *Target {
	viewmodel := new(Target)
	viewmodel.Target = t
	targetType := t.TargetType()
	viewmodel.templates = targetType.Templates()
	viewmodel.scripts = targetType.Scripts()

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

func (t *Target) Templates() map[string]string {
	return t.templates
}

func (t *Target) Scripts() map[string][]string {
	return t.scripts
}
