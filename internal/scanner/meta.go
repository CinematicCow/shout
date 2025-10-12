package scanner

import (
	"fmt"
	"os"
	"sort"
	"time"
)

func (s *Scanner) generateMeta(metaFile, name string, stats *Stats) error {
	file, err := os.Create(metaFile)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", cerr)
		}
	}()

	if _, err := fmt.Fprintf(file, "# %s - Meta Information\n\n", name); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	if _, err := fmt.Fprintf(file, "## Command\n```bash\n%s\n```\n\n", stats.Command); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	if _, err := fmt.Fprintf(file, "## Processed Files\n```\n%s```\n", s.generateMetaTree(stats.ProcessedFiles)); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if _, err := fmt.Fprintf(file, "## Statistics\n- Files processed: %d\n- Files skipped: %d\n- Generation time: %v\n\n",
		stats.FilesProcessed, stats.FilesSkipped, FormatDuration(stats.Duration)); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if _, err := fmt.Fprintf(file, "## File Extensions\n| Extension | Count |\n|-----------|-------|\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	exts := make([]string, 0, len(stats.FileStats))
	for ext := range stats.FileStats {
		exts = append(exts, ext)
	}
	sort.Strings(exts)
	for _, ext := range exts {
		if _, err := fmt.Fprintf(file, "| %-9s | %-5d |\n", ext, stats.FileStats[ext]); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	if len(stats.SkipPatterns) > 0 {
		if _, err := fmt.Fprintf(file, "\n## Skip Patterns Used\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		for _, pattern := range stats.SkipPatterns {
			if _, err := fmt.Fprintf(file, "- `%s`\n", pattern); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}
	return nil
}

func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return d.Round(time.Microsecond).String()
	} else if d < time.Second {
		return d.Round(time.Millisecond).String()
	}
	return d.Round(time.Second).String()
}
