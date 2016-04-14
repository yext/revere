package vm

type Component interface {
	Id() int64
}

type NamedComponent interface {
	Name() string
}

type DescriptiveComponent interface {
	Description() string
}

type SubprobeComponent interface {
	Subprobe() string
}
