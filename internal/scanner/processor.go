package scanner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (s *Scanner) processFiles(writer io.Writer, stats *Stats) error {
	absOutFile, _ := filepath.Abs(s.outFile)

	for _, dir := range s.Directories {
		filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil || shouldSkip(path, info, dir, absOutFile, s.SkipPatterns, s.Extensions, stats) {
				return err
			}

			content, _ := os.ReadFile(path)
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			fmt.Fprintf(writer, "\n## File: %s\n```%s\n%s\n```\n", path, ext, content)

			stats.FileStats[ext]++
			stats.FilesProcessed++
			stats.ProcessedFiles = append(stats.ProcessedFiles, path)
			return nil
		})
	}
	return nil
}

func shouldSkip(path string, info fs.FileInfo, dir, absOutFile string, skipPatterns, extensions []string, stats *Stats) bool {
	if info.IsDir() || strings.HasPrefix(info.Name(), ".") && path != dir {
		return true
	}
	if absPath, err := filepath.Abs(path); err == nil && absPath == absOutFile {
		return true
	}
	for _, pattern := range skipPatterns {
		if matched, _ := filepath.Match(pattern, info.Name()); matched {
			stats.FilesSkipped++
			return true
		}
	}
	if len(extensions) > 0 && !slices.Contains(extensions, strings.TrimPrefix(filepath.Ext(path), ".")) {
		stats.FilesSkipped++
		return true
	}
	return false
}
