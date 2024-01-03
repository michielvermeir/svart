package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	flags, args := InitializeCommandLine()

	Stderr.Printf("version %s\n", getAppVersion())

	if *flags.Version.Value {
		os.Exit(0)
	}

	os.Setenv("SVART_RELAXED_MODE", strconv.FormatBool(*flags.Relaxed.Value))

	if *flags.Prefix.Value != "" {
		os.Setenv("SVART_PREFIX", *flags.Prefix.Value)
	}

	if *flags.Filter.Value != "" {
		os.Setenv("SVART_DOTENV_FILTER", *flags.Filter.Value)
	}

	inputs := GetCommandLineInputs(flags)
	svartEnv := BuildExports(inputs)

	if len(svartEnv) == 0 {
		Stderr.Fatalf("no variables to re-export\n")
	}

	if len(args) < 1 {
		exportSvartEnv(svartEnv)
		return
	}

	executable := args[0]
	execArgs := args[1:]

	executablePath, err := exec.LookPath(executable)

	if err != nil {
		Stderr.Fatalf("something went wrong %s\n", err)
	}

	execEnv := append(os.Environ(), svartEnv...)
	command := exec.Command(executablePath, execArgs...)
	command.Env = execEnv

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	Stderr.Printf("executing %s\n", strings.Join(append([]string{executablePath}, execArgs...), " "))
	execErr := command.Run()

	if execErr != nil {
		Stderr.Fatalf("%s\n", execErr)
	}
}

func exportSvartEnv(svart []string) {
	for _, tfvar := range svart {
		fmt.Printf("export %s\n", tfvar)
	}
}
