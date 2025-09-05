package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CinematicCow/shout/internal/scanner"
	"github.com/CinematicCow/shout/internal/tui"
	"github.com/spf13/cobra"
)

var (
	version      = "1.0"
	outFile      = "llm.md"
	extensions   []string
	directories  []string
	skipPatterns []string
	interactive  bool
	meta         bool
	rootCmd      = &cobra.Command{
		Use:     "shout",
		Short:   "Project Dump for LLMs",
		Long:    `Shout v` + version + ` - Generates a single Markdown file of your project for LLMs`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && !interactive && len(extensions) == 0 && len(directories) == 0 && len(skipPatterns) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			if err := run(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.Flags().StringSliceVarP(&extensions, "extensions", "e", nil, "File extensions to include (comma-separated)")
	rootCmd.Flags().StringSliceVarP(&directories, "directories", "d", nil, "Directories to scan (comma-separated)")
	rootCmd.Flags().StringSliceVarP(&skipPatterns, "skip", "s", nil, "Patterns to skip (comma-separated)")
	rootCmd.Flags().StringVarP(&outFile, "output", "o", outFile, "Output file")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Use interactive TUI mode")
	rootCmd.Flags().BoolVarP(&meta, "meta", "m", false, "Generate meta information file")
}

func run(cmd *cobra.Command, args []string) error {
	if interactive {
		config, err := tui.StartTUI()
		if err != nil {
			return err
		}

		extensions = config.Extensions
		directories = config.Directories
		skipPatterns = config.SkipPatterns
		outFile = config.OutputFile
		meta = config.Meta
	}

	absOutFile, err := filepath.Abs(outFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output file: %v", err)
	}
	outFile = absOutFile

	if len(directories) == 0 {
		directories = []string{"."}
	}

	if _, err := os.Stat(".gitignore"); err == nil {
		content, err := os.ReadFile(".gitignore")
		if err == nil {
			lines := strings.SplitSeq(string(content), "\n")
			for line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					skipPatterns = append(skipPatterns, line)
				}
			}
		}
	}

	skipPatterns = append(skipPatterns, ".*")

	for _, dir := range directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not found", dir)
		}
	}

	s := scanner.New(extensions, directories, skipPatterns, outFile)
	projectName := filepath.Base(getCurrentDir())

	stats, err := s.Generate(outFile, projectName, meta)
	if err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", filepath.Base(outFile))
	fmt.Printf("Files: %d processed, %d skipped\n", stats.FilesProcessed, stats.FilesSkipped)
	fmt.Printf("Time: %v\n", scanner.FormatDuration(stats.Duration))

	if meta {
		fmt.Printf("Meta File: %s\n", stats.MetaFile)
	}

	return nil
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
