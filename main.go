package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"

	tfenv "tfvars/tfenv"

	"github.com/urfave/cli/v2"
)

func Map[T any](s []T, f func(T) T) []T {
	result := make([]T, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

func isStdinEmpty() bool {
	stdin, _ := os.Stdin.Stat()
	return stdin.Size() == 0
}

var (
	//go:embed VERSION
	res embed.FS
)

func getAppVersion() string {
	version, _ := res.ReadFile("VERSION")
	return string(version)
}

var app = &cli.App{
	Name:    "tfvars",
	Usage:   "re-exports the contents of dotenv files as TF_VAR_*",
	Version: getAppVersion(),
	Action: func(*cli.Context) error {
		if isStdinEmpty() {
			log.Fatal("tfvars: no input, stdin empty\n")
		}

		stdin := bufio.NewScanner(os.Stdin)
		tfvars := tfenv.Build(stdin)

		if len(tfvars) == 0 {
			fmt.Fprintf(os.Stderr, "tfvars: no variables in stdin to re-export\n")
			os.Exit(1)
		}

		for _, tfvar := range tfvars {
			fmt.Printf("export %s\n", tfvar)
		}

		return nil
	},
}

func main() {
	// Disable timestamp in log output
	log.SetFlags(0)

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "prints the version",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
