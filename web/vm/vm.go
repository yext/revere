// Package vm represents the intermediate step between a revere db struct and a
// renderable front-end struct.
package vm

type Component interface {
	Id() int64
}

type NamedComponent interface {
	Component
	ComponentName() string
}

type CreatableComponent interface {
	IsCreate() bool
}

type DeletableComponent interface {
	IsDelete() bool
}

func isCreate(c CreatableComponent) bool {
	return c.IsCreate()
}

func isDelete(c DeletableComponent) bool {
	return c.IsDelete()
}
