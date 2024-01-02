package main

import "flag"

type arg[T string | bool] struct {
	Name  string
	Value *T
	IsSet func() bool
}

type Args struct {
	File          arg[string]
	FromEnv       arg[bool]
	FromStdin     arg[bool]
	Version       arg[bool]
	AllowlistFile arg[string]
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

func InitializeCommandLine() *Args {
	args := &Args{
		File:          define[string]("from-file", flag.String("from-file", "", "read from file")),
		FromEnv:       define[bool]("from-env", flag.Bool("from-env", true, "read from environment variables")),
		FromStdin:     define[bool]("from-stdin", flag.Bool("from-stdin", false, "read from stdin")),
		Version:       define[bool]("version", flag.Bool("version", false, "print version")),
		AllowlistFile: define[string]("allowlist-file", flag.String("allowlist", "", "only re-export allowlisted variables")),
	}

	// Must parse after all flags are defined
	flag.Parse()
	return args
}
