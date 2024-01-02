package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"svart/cli"
	"svart/exports"
	"svart/inputs"
)

var (
	//go:embed VERSION
	res embed.FS
)

var Stderr = log.New(os.Stderr, "svart: ", 0)

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

	if *args.Prudent.Value {
		os.Setenv("SVART_BLOCKED", "*")
	}

	inputs := inputs.Get(args)
	svart := exports.Build(inputs)

	if len(svart) == 0 {
		Stderr.Fatalf("no variables to re-export\n")
	}

	for _, tfvar := range svart {
		fmt.Printf("export %s\n", tfvar)
	}

	// if err != nil {
	// 	stderr.Fatalf("something went wrong %s\n", err)
	// }
}
