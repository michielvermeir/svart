package main

import (
	"os"
	"strconv"
	"strings"

	filepath "path/filepath"
)

func Getenv(name, fallback string) string {
	if _, found := os.LookupEnv(name); !found {
		return fallback
	}
	return os.Getenv(name)
}

func GetAllowed() []string {
	allowlistFile := os.Getenv("SVART_DOTENV_FILTER")

	// Stderr.Printf("allowlist file %s\n", file)
	if len(allowlistFile) != 0 {
		allowed := GetEnvPairNamesFromFile(allowlistFile)
		// Stderr.Printf("allowed names %v\n", allowed)
		return allowed
	}

	allowlistPatterns := os.Getenv("SVART_FILTERS")
	if len(allowlistPatterns) != 0 {
		return strings.Split(allowlistPatterns, ",")
	}

	if relaxed, _ := strconv.ParseBool(os.Getenv("SVART_RELAXED_MODE")); relaxed {
		return []string{"*"}
	}

	// Strict mode by default
	return []string{}
}

func IsExportAllowed(name string) bool {
	allowed := GetAllowed()

	for _, pattern := range allowed {
		if match, _ := filepath.Match(pattern, name); match {
			return true
		}
	}

	return false
}

func IsBlocked(name string) bool {
	blocked := Getenv("SVART_BLOCKED", "AWS_*")

	// Stderr.Printf("blocked %s\n", blocked)

	for _, pattern := range strings.Split(blocked, ",") {
		if len(pattern) == 0 {
			continue
		}

		match, _ := filepath.Match(pattern, name)
		return match
	}

	return false
}
