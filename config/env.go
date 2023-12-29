package config

import (
	"log"
	"os"
	"strings"

	filepath "path/filepath"
)

var Stderr = log.New(os.Stderr, "tfvars: ", 0)

func Getenv(name, fallback string) string {
	if _, found := os.LookupEnv(name); !found {
		return fallback
	}
	return os.Getenv(name)
}

func IsAllowed(name string) bool {
	allowed := Getenv("TFVARS_ALLOWED", "*")

	// Stderr.Printf("allowed %s\n", allowed)

	for _, pattern := range strings.Split(allowed, ",") {
		match, _ := filepath.Match(pattern, name)
		return match
	}

	return false
}

func IsBlocked(name string) bool {
	blocked := Getenv("TFVARS_BLOCKED", "AWS_*")

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
