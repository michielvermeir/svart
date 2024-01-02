package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/direnv/go-dotenv"
)

var SINGLE_LINE = regexp.MustCompile(
	regexp.MustCompile("\\s+").ReplaceAllLiteralString(
		regexp.MustCompile("\\s+# .*").ReplaceAllLiteralString(dotenv.LINE, ""), ""))

type Pair struct {
	Name  string
	Value string
}

func (m *Pair) AsPrefixed() string {
	prefix := Getenv("SVART_PREFIX", "TF_VAR_")
	return fmt.Sprintf("%s%s=%s", prefix, strings.ToLower(m.Name), m.Value)
}

func normalize(fileOrPath string) string {
	if filepath.IsAbs(fileOrPath) {
		return fileOrPath
	}

	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, fileOrPath)
}

func GetEnvPairsFromFile(fileOrPath string) []Pair {
	filePath := normalize(fileOrPath)

	// Stderr.Printf("reading from file %s\n", filePath)
	file, err := os.Open(filePath)

	if err != nil {
		Stderr.Fatalf("%s\n", err)
	}

	buffer := bufio.NewScanner(file)

	var pairs = []Pair{}

	for buffer.Scan() {
		currentLine := buffer.Text()
		// Stderr.Printf("current line %v\n", currentLine)

		match := parseLine(currentLine)
		if match != nil {
			pairs = append(pairs, *match)
		}
	}

	return pairs
}

func GetEnvPairNamesFromFile(file string) []string {
	pairs := GetEnvPairsFromFile(file)
	// Stderr.Printf("pairs in file %v\n", pairs)

	var names = []string{}
	for _, pair := range pairs {
		names = append(names, pair.Name)
	}

	return names
}

func GetEnvPairValuesFromFile(file string) []string {
	pairs := GetEnvPairsFromFile(file)

	var values = []string{}

	for _, pair := range pairs {
		values = append(values, pair.Value)
	}

	return values
}

func parseLine(line string) *Pair {
	match := SINGLE_LINE.MatchString(line)

	if match {
		parsed, _ := dotenv.Parse(line)
		for name, value := range parsed {
			return &Pair{Name: name, Value: value}
		}
	}

	return nil
}

func BuildExports(buffer *bufio.Scanner) []string {
	var svart = []string{}

	for buffer.Scan() {
		currentLine := buffer.Text()
		match := parseLine(currentLine)

		if match != nil && IsExportAllowed(match.Name) {
			svart = append(svart, match.AsPrefixed())
		}
	}

	if err := buffer.Err(); err != nil {
		fmt.Println(err)
	}

	sort.Strings(svart)
	return svart
}
