package inputs

import (
	"bufio"
	"log"
	"os"
	"strings"
	"svart/cli"
)

var Stderr = log.New(os.Stderr, "svart: ", 0)

func isStdinEmpty() bool {
	stdin, _ := os.Stdin.Stat()
	return stdin.Size() == 0
}

type input interface {
	get(args *cli.Args) *bufio.Scanner
}

type stdin struct{}

func (input *stdin) get(args *cli.Args) *bufio.Scanner {
	Stderr.Print("reading from stdin\n")

	if isStdinEmpty() {
		Stderr.Fatal("stdin empty, no input\n")
	}

	return bufio.NewScanner(os.Stdin)
}

type file struct{}

func (input *file) get(args *cli.Args) *bufio.Scanner {
	Stderr.Printf("reading from file %s\n", *args.File.Value)

	file, err := os.Open(*args.File.Value)
	if err != nil {
		Stderr.Fatalf("%s\n", err)
	}

	return bufio.NewScanner(file)
}

type environ struct{}

func (input *environ) get(args *cli.Args) *bufio.Scanner {
	Stderr.Print("reading from environment variables\n")

	return bufio.NewScanner(
		strings.NewReader(strings.Join(os.Environ(), "\n")),
	)
}

func reify(args *cli.Args) input {
	if *args.FromStdin.Value {
		return &stdin{}
	}

	if args.File.IsSet() {
		return &file{}
	}

	// fallback to environment variables
	return &environ{}
}

// Gets a scanner for the appropriate input source (env, file or stdin)
func Get(cliArgs *cli.Args) *bufio.Scanner {
	return reify(cliArgs).get(cliArgs)
}
