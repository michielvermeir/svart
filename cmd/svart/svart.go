package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
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

	if *args.AllowlistFile.Value != "" {
		os.Setenv("SVART_ALLOWLIST_FILE", *args.AllowlistFile.Value)
	}

	inputs := GetCommandLineInputs(args)
	svartEnv := BuildExports(inputs)

	if len(svartEnv) == 0 {
		Stderr.Fatalf("no variables to re-export\n")
	}

	if len(os.Args) < 2 {
		exportSvartEnv(svartEnv)
	}

	binary := os.Args[1]
	execArgs := os.Args[1:]

	binaryPath, err := exec.LookPath(binary)

	if err != nil {
		Stderr.Fatalf("something went wrong %s\n", err)
	}

	execEnv := append(os.Environ(), svartEnv...)
	execErr := syscall.Exec(binaryPath, execArgs, execEnv)

	if execErr != nil {
		Stderr.Fatalf("%s\n", execErr)
	}
}

func exportSvartEnv(svart []string) {
	for _, tfvar := range svart {
		fmt.Printf("export %s\n", tfvar)
	}
}
