package conf

var options *Options

type Option func(opt *Options)
type Options struct {
	DataPath string
}

func DefOption() *Options {
	return &Options{
		DataPath: "./",
	}
}

func LoadOption(opts ...Option) *Options {
	options = DefOption()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func GetOptions() *Options {
	return options
}

func WithDataPath(p string) Option {
	return func(opt *Options) {
		opt.DataPath = p
	}
}
