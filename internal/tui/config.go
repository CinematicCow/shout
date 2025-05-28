package tui

import "strings"

type Config struct {
	Extensions   []string
	Directories  []string
	SkipPatterns []string
	OutputFile   string
	Meta         bool
}

func splitCommaString(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
