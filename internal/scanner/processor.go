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
	cwd, _ := os.Getwd()
	absOutFile, _ := filepath.Abs(s.outFile)
	processedFiles := make(map[string]bool)

	for _, path := range s.Directories {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		info, err := os.Stat(absPath)
		if err != nil {
			continue
		}

		if info.IsDir() {
			if err := filepath.Walk(absPath, func(walkPath string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				return s.processFile(walkPath, info, writer, stats, cwd, absOutFile, processedFiles)
			}); err != nil {
				return fmt.Errorf("failed to walk directory %s: %w", absPath, err)
			}
		} else {
			if err := s.processFile(absPath, info, writer, stats, cwd, absOutFile, processedFiles); err != nil {
				return fmt.Errorf("failed to process file %s: %w", absPath, err)
			}
		}
	}
	return nil
}

func (s *Scanner) processFile(path string, info fs.FileInfo, writer io.Writer, stats *Stats, cwd, absOutFile string, processedFiles map[string]bool) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if processedFiles[absPath] {
		return nil
	}
	processedFiles[absPath] = true

	if info.IsDir() {
		return nil
	}

	if s.shouldSkip(path, absPath, cwd, absOutFile, stats) {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fileTokens := EstimateTokens(content)
	stats.TotalTokens += fileTokens

	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	if _, err := fmt.Fprintf(writer, "\n## File: %s\n```%s\n%s\n```\n", path, ext, content); err != nil {
		return fmt.Errorf("failed to write file content to writer: %w", err)
	}

	stats.FilesProcessed++
	stats.ProcessedFiles = append(stats.ProcessedFiles, path)
	stats.FileStats[ext]++

	return nil
}

func (s *Scanner) shouldSkip(path, absPath, cwd, absOutFile string, stats *Stats) bool {
	base := filepath.Base(path)

	if strings.HasPrefix(base, ".") {
		return true
	}

	if absPath == absOutFile {
		return true
	}

	for _, pattern := range s.SkipPatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			stats.FilesSkipped++
			return true
		}

		if relPath, err := filepath.Rel(cwd, absPath); err == nil {
			if matched, _ := filepath.Match(pattern, relPath); matched {
				stats.FilesSkipped++
				return true
			}
		}

		if matched, _ := filepath.Match(pattern, absPath); matched {
			stats.FilesSkipped++
			return true
		}

		// Convert pattern to absolute path for proper directory matching
		absPattern, err := filepath.Abs(pattern)
		if err == nil {
			// Ensure pattern ends with separator for directory matching
			if !strings.HasSuffix(absPattern, string(filepath.Separator)) {
				absPattern += string(filepath.Separator)
			}
			if strings.HasPrefix(absPath, absPattern) {
				stats.FilesSkipped++
				return true
			}
		}
	}

	if len(s.Extensions) > 0 {
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if !slices.Contains(s.Extensions, ext) {
			stats.FilesSkipped++
			return true
		}
	}

	return false
}
