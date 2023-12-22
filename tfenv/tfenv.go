package tfenv

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/direnv/go-dotenv"
)

var SINGLE_LINE = regexp.MustCompile(
	regexp.MustCompile("\\s+").ReplaceAllLiteralString(
		regexp.MustCompile("\\s+# .*").ReplaceAllLiteralString(dotenv.LINE, ""), ""))

type Line struct {
	Name  string
	Value string
}

func (m *Line) AsTfvar() string {
	return fmt.Sprintf("TF_VAR_%s=%s", strings.ToLower(m.Name), m.Value)
}

func matchLine(line string) (*Line, error) {
	match := SINGLE_LINE.MatchString(line)

	if match {
		parsed, _ := dotenv.Parse(line)
		for key, value := range parsed {
			return &Line{Name: key, Value: value}, nil
		}
	}

	return nil, nil
}

func Build(buffer *bufio.Scanner) []string {
	var tfvars = []string{}

	for buffer.Scan() {
		currentLine := buffer.Text()
		match, _ := matchLine(currentLine)

		if match != nil {
			tfvars = append(tfvars, match.AsTfvar())
		}
	}

	if err := buffer.Err(); err != nil {
		fmt.Println(err)
	}

	return tfvars
}
