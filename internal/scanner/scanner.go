package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Scanner struct {
	Extensions   []string
	Directories  []string
	SkipPatterns []string
	outFile      string
}

func New(extensions, dirs, skip []string, outFile string) *Scanner {
	return &Scanner{
		Extensions:   extensions,
		Directories:  dirs,
		SkipPatterns: skip,
		outFile:      outFile,
	}
}

func confirmOverwrite(filePath string, force bool) (bool, error) {
	if force {
		return true, nil
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return true, nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("File '%s' exists. Overwrite? [y/N]: ", filepath.Base(filePath))
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}
