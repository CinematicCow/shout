package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func EstimateTokens(content []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	tc := 0

	for scanner.Scan() {
		line := scanner.Text()
		if isCodeLine(line) {
			tc += estimateCodeTokens(line)
		} else {
			tc += len(line) / 4
		}
	}
	return tc
}

func isCodeLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return false
	}

	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
		return false
	}

	hasBraces := strings.Contains(line, "{") || strings.Contains(line, "}")
	hasParens := strings.Contains(line, "(") || strings.Contains(line, ")")
	hasOperators := strings.ContainsAny(line, "=+-*/%<>!&|^~")

	return hasBraces || hasParens || hasOperators
}

func estimateCodeTokens(line string) int {
	tokens := 0
	ct := false

	for _, r := range line {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			if !ct {
				tokens++
				ct = true
			}
		} else if unicode.IsSpace(r) {
			ct = false
		} else {
			if !unicode.IsSpace(r) {
				tokens++
				ct = false
			}
		}
	}

	if tokens == 0 && len(strings.TrimSpace(line)) > 0 {
		tokens = 1
	}

	return tokens
}

func CountFileTokens(filepath string) (int, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	return EstimateTokens(content), nil
}
