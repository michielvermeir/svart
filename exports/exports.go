package exports

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"svart/config"

	"github.com/direnv/go-dotenv"
)

var SINGLE_LINE = regexp.MustCompile(
	regexp.MustCompile("\\s+").ReplaceAllLiteralString(
		regexp.MustCompile("\\s+# .*").ReplaceAllLiteralString(dotenv.LINE, ""), ""))

type Line struct {
	Name  string
	Value string
}

func (m *Line) AsPrefixed() string {
	prefix := config.Getenv("SVART_PREFIX", "TF_VAR_")
	return fmt.Sprintf("%s%s=%s", prefix, strings.ToLower(m.Name), m.Value)
}

func matchLine(line string) *Line {
	match := SINGLE_LINE.MatchString(line)

	if match {
		parsed, _ := dotenv.Parse(line)
		for name, value := range parsed {
			if config.IsAllowed(name) || !config.IsBlocked(name) {
				return &Line{Name: name, Value: value}
			}
		}
	}

	return nil
}

func Build(buffer *bufio.Scanner) []string {
	var svart = []string{}

	for buffer.Scan() {
		currentLine := buffer.Text()
		match := matchLine(currentLine)

		if match != nil {
			svart = append(svart, match.AsPrefixed())
		}
	}

	if err := buffer.Err(); err != nil {
		fmt.Println(err)
	}

	sort.Strings(svart)
	return svart
}
