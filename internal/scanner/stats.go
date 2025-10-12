package scanner

import (
	"os"
	"strings"
	"time"
)

type Stats struct {
	FilesProcessed int
	FilesSkipped   int
	Duration       time.Duration
	MetaFile       string
	FileStats      map[string]int
	ProcessedFiles []string
	SkipPatterns   []string
	Command        string
	TotalTokens    int
}

func BuildCommand() string {
	return strings.Join(os.Args, " ")
}
