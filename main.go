package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tfenv "tfvars/tfenv"
)

func Map[T any](s []T, f func(T) T) []T {
	result := make([]T, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	tfvars := tfenv.Build(stdin)

	fmt.Fprintf(os.Stderr, "tfvars: export %s\n",
		strings.Join(Map(tfvars, func(tfvar string) string {
			return fmt.Sprintf("+%s", strings.Split(tfvar, "=")[0])
		}), " "))

	for _, tfvar := range tfvars {
		fmt.Printf("export %s\n", tfvar)
	}
}
