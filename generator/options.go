package generator

// Options allows to customize the code generation process.
type Options struct {
	ModuleName   string
	PackageName  string
	ExtraImports []ImportSpec
	// Map of tuple definitions to existing struct names,
	// to avoid generating duplicate structs
	ExternalTuples map[string]string
	Prefix         string
	Stdlib         bool
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		PackageName:    "abi",
		ExtraImports:   []ImportSpec{},
		ExternalTuples: make(map[string]string),
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

func ModuleName(name string) Option {
	return func(o *Options) {
		o.ModuleName = name
	}
}

func Prefix(p string) Option {
	return func(o *Options) {
		o.Prefix = p
	}
}

func Stdlib(s bool) Option {
	return func(o *Options) {
		o.Stdlib = s
	}
}

func ExtraImports(imports []ImportSpec) Option {
	return func(o *Options) {
		o.ExtraImports = imports
	}
}

func ExternalTuples(m map[string]string) Option {
	return func(o *Options) {
		o.ExternalTuples = m
	}
}
