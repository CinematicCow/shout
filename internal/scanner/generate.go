package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Scanner) Generate(outFile, name string, meta bool) (*Stats, error) {
	start := time.Now()
	file, _ := os.Create(outFile)
	defer file.Close()

	fmt.Fprintf(file, "# Project: %s\n\n", name)

	stats := &Stats{
		FileStats:      make(map[string]int),
		ProcessedFiles: []string{},
		SkipPatterns:   s.SkipPatterns,
		Command:        BuildCommand(),
	}

	tree, _ := s.generateTree(s.Directories)
	fmt.Fprintf(file, "## Project Structure\n```\n%s\n```\n", tree)

	s.processFiles(file, stats)
	stats.Duration = time.Since(start)

	if meta {
		metaFile := strings.TrimSuffix(outFile, filepath.Ext(outFile)) + ".meta.md"
		stats.MetaFile = metaFile
		s.generateMeta(metaFile, name, stats)
	}
	return stats, nil
}
