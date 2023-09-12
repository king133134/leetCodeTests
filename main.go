package leetcodeTests

type Options struct {
	mod  string
	port int
	dir  string
}

type Option func(options *Options)

func Mod(mod string) Option {
	return func(options *Options) {
		options.mod = mod
	}
}

func Port(port int) Option {
	return func(options *Options) {
		options.port = port
	}
}

func Dir(dir string) Option {
	return func(options *Options) {
		options.dir = dir
	}
}

func Run(opts ...Option) {
	opt := &Options{
		mod:  "http",
		port: 8080,
		dir:  ".",
	}

	for _, o := range opts {
		o(opt)
	}
	switch opt.mod {
	case "http":
		HttpStart(opt.port)
	case "file":
		FileRun(opt.dir)
	case "input":
		InputRun()
	default:
		panic("mod option is error!")
	}
	return
}
