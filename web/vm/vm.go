package vm

type Component interface {
	Id() int64
}

type CreatableComponent interface {
	Create() bool
}

type DeletableComponent interface {
	Del() bool
}

func isCreate(c CreatableComponent) bool {
	return c.Create()
}

func isDelete(c DeletableComponent) bool {
	return c.Del()
}
