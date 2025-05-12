package scanner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Scanner struct {
	Extensions   []string
	Directories  []string
	SkipPatterns []string
}

func New(extenstions, dirs, skip []string) *Scanner {
	return &Scanner{
		Extensions:   extenstions,
		Directories:  dirs,
		SkipPatterns: skip,
	}
}

func (s *Scanner) Generate(outFile, name string) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# Project: %s\n\n", name)

	tree, err := s.generateTree()
	if err != nil {
		return err
	}

	fmt.Fprintf(file, "## Project Structure\n```\n%s\n```\n", tree)

	err = s.processFiles(file)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scanner) generateTree() (string, error) {
	treeCmd := exec.Command("tree", append([]string{"--noreport"}, s.Directories...)...)
	output, err := treeCmd.Output()
	if err == nil {
		return string(output), nil
	}

	var sb strings.Builder
	for _, dir := range s.Directories {
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				relPath, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}

				if relPath == "." {
					sb.WriteString(filepath.Base(dir) + "\n")
				} else {
					depth := strings.Count(relPath, string(os.PathSeparator))
					indent := strings.Repeat("| ", depth) + "|- "
					sb.WriteString(indent + filepath.Base(path) + "\n")
				}
			}
			return nil
		})

		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func (s *Scanner) processFiles(writer io.Writer) error {
	for _, dir := range s.Directories {
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			for _, pattern := range s.SkipPatterns {
				if matched, err := filepath.Match(pattern, filepath.Base(path)); err == nil && matched {
					return nil
				}
			}

			if len(s.Extensions) > 0 {
				ext := strings.TrimPrefix(filepath.Ext(path), ".")
				matched := false
				for _, allowedExt := range s.Extensions {
					if ext == allowedExt {
						matched = true
						break
					}
				}

				if !matched {
					return nil
				}
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			fmt.Fprintf(writer, "\n## File: %s\n```%s\n%s\n```\n", path, ext, string(content))

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
