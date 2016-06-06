package vm

type Component interface {
	Id() int64
}

type DeletableComponent interface {
	Del() bool
}

func isCreate(c Component) bool {
	return c.Id() == 0
}

func isDelete(c DeletableComponent) bool {
	return c.Del()
}
