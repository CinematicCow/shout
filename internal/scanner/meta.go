package scanner

import (
	"fmt"
	"os"
	"sort"
	"time"
)

func (s *Scanner) generateMeta(metaFile, name string, stats *Stats) error {
	file, _ := os.Create(metaFile)
	defer file.Close()

	fmt.Fprintf(file, "# %s - Meta Information\n\n", name)
	fmt.Fprintf(file, "## Command\n```bash\n%s\n```\n\n", stats.Command)
	fmt.Fprintf(file, "## Processed Files\n```\n%s```\n", s.generateMetaTree(stats.ProcessedFiles))

	fmt.Fprintf(file, "## Statistics\n- Files processed: %d\n- Files skipped: %d\n- Generation time: %v\n\n",
		stats.FilesProcessed, stats.FilesSkipped, FormatDuration(stats.Duration))

	fmt.Fprintf(file, "## File Extensions\n| Extension | Count |\n|-----------|-------|\n")
	exts := make([]string, 0, len(stats.FileStats))
	for ext := range stats.FileStats {
		exts = append(exts, ext)
	}
	sort.Strings(exts)
	for _, ext := range exts {
		fmt.Fprintf(file, "| %-9s | %-5d |\n", ext, stats.FileStats[ext])
	}

	if len(stats.SkipPatterns) > 0 {
		fmt.Fprintf(file, "\n## Skip Patterns Used\n")
		for _, pattern := range stats.SkipPatterns {
			fmt.Fprintf(file, "- `%s`\n", pattern)
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
