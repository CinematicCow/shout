package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CinematicCow/shout/internal/scanner"
	"github.com/CinematicCow/shout/internal/tui"
	"github.com/spf13/cobra"
)

var (
	version      = "1.0.0"
	outFile      = "llm.md"
	extensions   []string
	directories  []string
	skipPatterns []string
	interactive  bool
	rootCmd      = &cobra.Command{
		Use:     "shout",
		Short:   "Project Dump for LLM analysis",
		Long:    `Shout v` + version + ` - Generates a single Markdown file of your project for LLM analysis`,
		Version: version,
		RunE:    run,
	}
)

func init() {
	rootCmd.Flags().StringSliceVarP(&extensions, "extensions", "e", nil, "File extensions to include (comma-separated)")
	rootCmd.Flags().StringSliceVarP(&directories, "directories", "d", nil, "Directories to scan (comma-separated)")
	rootCmd.Flags().StringSliceVarP(&skipPatterns, "skip", "s", nil, "Patterns to skip (comma-separated)")
	rootCmd.Flags().StringVarP(&outFile, "output", "o", outFile, "Output file")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Use interactive TUI mode")
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
	}

	if len(directories) == 0 {
		directories = []string{"."}
	}

	for _, dir := range directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not found", dir)
		}
	}

	fmt.Println("Generating project dump...")

	s := scanner.New(extensions, directories, skipPatterns)
	projectName := filepath.Base(getCurrentDir())
	if err := s.Generate(outFile, projectName); err != nil {
		return err
	}
	fmt.Printf("Generated: %s\n", outFile)
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
