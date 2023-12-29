package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"tfvars/cli"
	"tfvars/exports"
	"tfvars/inputs"
)

var (
	//go:embed VERSION
	res embed.FS
)

var Stderr = log.New(os.Stderr, "tfvars: ", 0)

func getAppVersion() string {
	version, _ := res.ReadFile("VERSION")
	return string(version)
}

type App struct{}

func main() {
	// Disable timestamp in log output
	log.SetFlags(0)

	args := cli.Initialize()

	Stderr.Printf("version %s\n", getAppVersion())

	if *args.Version.Value {
		os.Exit(0)
	}

	inputs := inputs.Get(args)
	tfvars := exports.Build(inputs)

	if len(tfvars) == 0 {
		Stderr.Fatalf("no variables to re-export\n")
	}

	for _, tfvar := range tfvars {
		fmt.Printf("export %s\n", tfvar)
	}

	// if err != nil {
	// 	stderr.Fatalf("something went wrong %s\n", err)
	// }
}
