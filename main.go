package main

import (
	"embed"
	"fmt"
	"log"
	"os"
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

	if *args.AllowlistFile.Value != "" {
		os.Setenv("SVART_ALLOWLIST_FILE", *args.AllowlistFile.Value)
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
