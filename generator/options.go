package generator

// Options allows to customize the code generation process.
type Options struct {
	PackageName  string
	ExtraImports []string
	// Map of tuple definitions to existing struct names,
	// to avoid generating duplicate structs
	TupleMap map[string]string
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		PackageName:  "abi",
		ExtraImports: []string{},
		TupleMap:     make(map[string]string),
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type Option func(*Options)

func PackageName(name string) Option {
	return func(o *Options) {
		o.PackageName = name
	}
}

func ExtraImports(imports []string) Option {
	return func(o *Options) {
		o.ExtraImports = imports
	}
}

func TupleMap(m map[string]string) Option {
	return func(o *Options) {
		o.TupleMap = m
	}
}
