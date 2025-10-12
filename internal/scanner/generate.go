package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Scanner) Generate(outFile, name string, meta bool, git bool, gitLimit int) (*Stats, error) {
	start := time.Now()

	file, err := os.Create(outFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", cerr)
		}
	}()

	if _, err := fmt.Fprintf(file, "# Project: %s\n\n", name); err != nil {
		return nil, fmt.Errorf("failed to write to file: %w", err)
	}

	stats := &Stats{
		FileStats:      make(map[string]int),
		ProcessedFiles: []string{},
		SkipPatterns:   s.SkipPatterns,
		Command:        BuildCommand(),
	}

	tree, _ := s.generateTree(s.Directories)

	if _, err := fmt.Fprintf(file, "## Project Structure\n```\n%s\n```\n", tree); err != nil {
		return nil, fmt.Errorf("failed to write to file: %w", err)
	}

	if git {
		commits, err := GetGit(gitLimit)
		if err == nil {
			if _, err := fmt.Fprintf(file, "%s\n", FormatGit(commits)); err != nil {
				return nil, fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}

	if err := s.processFiles(file, stats); err != nil {
		return nil, fmt.Errorf("failed to process files: %w", err)
	}
	stats.Duration = time.Since(start)

	if meta {
		metaFile := strings.TrimSuffix(filepath.Base(outFile), filepath.Ext(outFile)) + ".meta.md"
		stats.MetaFile = metaFile
		metaPath := filepath.Join(filepath.Dir(outFile), metaFile)
		if err := s.generateMeta(metaPath, name, stats); err != nil {
			return nil, fmt.Errorf("failed to generate meta file: %w", err)
		}
	}
	return stats, nil
}
