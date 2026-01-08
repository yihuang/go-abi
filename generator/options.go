package generator

// Options allows to customize the code generation process.
type Options struct {
	PackageName  string
	ExtraImports []ImportSpec
	// Map of tuple definitions to existing struct names,
	// to avoid generating duplicate structs
	ExternalTuples map[string]string
	Prefix         string
	Stdlib         bool
	UseUint256     bool   // Use holiman/uint256 for uint256 types instead of *big.Int
	BuildTag       string // Build tag to add to generated file (e.g., "uint256")
	GenerateLazy   bool   // Generate lazy decoding View types
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

func UseUint256(use bool) Option {
	return func(o *Options) {
		o.UseUint256 = use
	}
}

func BuildTag(tag string) Option {
	return func(o *Options) {
		o.BuildTag = tag
	}
}

func GenerateLazy(enable bool) Option {
	return func(o *Options) {
		o.GenerateLazy = enable
	}
}
