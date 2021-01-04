package options

type Option func(d *Demo)

type Demo struct {
	Name string
	Addr string
	Num  string
}

func NewDemo(options ...Option) *Demo {
	d := &Demo{}
	for _, option := range options {
		option(d)
	}
	return d
}

func WithName(name string) Option {
	return func(d *Demo) {
		d.Name = name
	}
}
func WithAddr(addr string) Option {
	return func(d *Demo) {
		d.Name = addr
	}
}
