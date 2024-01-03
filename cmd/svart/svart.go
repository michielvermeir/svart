package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
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

	args := InitializeCommandLine()

	Stderr.Printf("version %s\n", getAppVersion())

	if *args.Version.Value {
		os.Exit(0)
	}

	os.Setenv("SVART_RELAXED_MODE", strconv.FormatBool(*args.Relaxed.Value))

	if *args.Prefix.Value != "" {
		os.Setenv("SVART_PREFIX", *args.Prefix.Value)
	}

	if *args.Filter.Value != "" {
		os.Setenv("SVART_ALLOWLIST_FILE", *args.Filter.Value)
	}

	inputs := GetCommandLineInputs(args)
	svart := BuildExports(inputs)

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
