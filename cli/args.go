package cli

import "flag"

type arg[T string | bool] struct {
	Name  string
	Value *T
	IsSet func() bool
}

type Args struct {
	Version   arg[bool]
	File      arg[string]
	FromStdin arg[bool]
	FromEnv   arg[bool]
	Prudent   arg[bool]
}

func define[T string | bool](name string, value *T) arg[T] {
	version := flag.Lookup(name)

	return arg[T]{
		Name:  name,
		Value: value,
		IsSet: func() bool {
			return version.DefValue != version.Value.String()
		},
	}
}

func Initialize() *Args {
	args := &Args{
		Version:   define[bool]("version", flag.Bool("version", false, "print version")),
		FromStdin: define[bool]("from-stdin", flag.Bool("from-stdin", false, "read from stdin")),
		FromEnv:   define[bool]("from-env", flag.Bool("from-env", true, "read from environment variables")),
		File:      define[string]("from-file", flag.String("from-file", "", "read from file")),
		Prudent:   define[bool]("prudent", flag.Bool("prudent", false, "prudent mode (all variables blocked)")),
	}

	// Must parse after all flags are defined
	flag.Parse()
	return args
}
