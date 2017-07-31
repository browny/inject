package example

type Logger interface {
	Log(format string, a ...interface{})
}

type Food interface {
	GetRice()
}

type Machine interface {
	Run(n int) error
}

type Transport interface {
	Fly(src, dst string)
}
