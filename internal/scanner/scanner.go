package scanner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type Scanner struct {
	Extensions   []string
	Directories  []string
	SkipPatterns []string
	outFile      string
}

type Stats struct {
	FilesProcessed int
	FilesSkipped   int
	Duration       time.Duration
	MetaFile       string
	FileStats      map[string]int
	ProcessedFiles []string
	SkipPatterns   []string
	Command        string
}

func New(extenstions, dirs, skip []string, outFile string) *Scanner {
	return &Scanner{
		Extensions:   extenstions,
		Directories:  dirs,
		SkipPatterns: skip,
		outFile:      outFile,
	}
}

func (s *Scanner) Generate(outFile, name string, meta bool) (*Stats, error) {
	start := time.Now()
	file, err := os.Create(outFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Fprintf(file, "# Project: %s\n\n", name)

	stats := &Stats{
		FileStats:      make(map[string]int),
		ProcessedFiles: []string{},
		SkipPatterns:   s.SkipPatterns,
		Command:        buildCommand(),
	}

	tree, err := s.generateTree(s.Directories)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(file, "## Project Structure\n```\n%s\n```\n", tree)

	err = s.processFiles(file, stats)
	if err != nil {
		return nil, err
	}

	stats.Duration = time.Since(start)

	if meta {
		metaFile := strings.TrimSuffix(outFile, filepath.Ext(outFile)) + ".meta.md"
		stats.MetaFile = metaFile
		err := s.generateMeta(metaFile, name, stats)
		if err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func (s *Scanner) generateTree(files []string) (string, error) {
	treeCmd := exec.Command("tree", append([]string{"--noreport"}, files...)...)
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

func (s *Scanner) processFiles(writer io.Writer, stats *Stats) error {

	absOutFile, err := filepath.Abs(s.outFile)
	if err != nil {
		return err
	}

	for _, dir := range s.Directories {
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(filepath.Base(path), ".") && path != dir {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			absPath, err := filepath.Abs(path)
			if err == nil && absPath == absOutFile {
				return nil
			}

			if info.IsDir() {
				return nil
			}
			for _, pattern := range s.SkipPatterns {
				if matched, err := filepath.Match(pattern, filepath.Base(path)); err == nil && matched {
					stats.FilesSkipped++
					return nil
				}
			}
			if len(s.Extensions) > 0 {
				ext := strings.TrimPrefix(filepath.Ext(path), ".")
				matched := slices.Contains(s.Extensions, ext)
				if !matched {
					stats.FilesSkipped++
					return nil
				}
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			fmt.Fprintf(writer, "\n## File: %s\n```%s\n%s\n```\n", path, ext, string(content))

			stats.FilesProcessed++
			stats.ProcessedFiles = append(stats.ProcessedFiles, path)
			stats.FileStats[ext]++

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scanner) generateMeta(metaFile, name string, stats *Stats) error {
	file, err := os.Create(metaFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# %s - Meta Information\n\n", name)

	fmt.Fprintf(file, "## Command\n```bash\n%s\n```\n\n", stats.Command)

	fmt.Fprintf(file, "## Processed Files\n```\n")
	tree := s.generateMetaTree(stats.ProcessedFiles)
	fmt.Fprintf(file, "%s```\n", tree)

	fmt.Fprintf(file, "## Statistics\n")
	fmt.Fprintf(file, "- **Files processed**: %d\n", stats.FilesProcessed)
	fmt.Fprintf(file, "- **Files skipped**: %d\n", stats.FilesSkipped)
	fmt.Fprintf(file, "- **Generation time**: %v\n\n", stats.Duration)

	fmt.Fprintf(file, "## File Extensions\n\n")
	fmt.Fprintf(file, "| Extension | Count |\n")
	fmt.Fprintf(file, "|-----------|-------|\n")
	for ext, count := range stats.FileStats {
		extDisplay := ext
		if ext == "" {
			extDisplay = "(no extension)"
		}
		fmt.Fprintf(file, "|     %s    |   %d   |\n\n", extDisplay, count)
	}

	if len(stats.SkipPatterns) > 0 {
		fmt.Fprintf(file, "## Skip Patterns Used\n")
		for _, pattern := range stats.SkipPatterns {
			fmt.Fprintf(file, "- `%s`\n", pattern)
		}
		fmt.Fprintf(file, "\n")
	}

	return nil
}

func (s *Scanner) generateMetaTree(files []string) string {
	if len(files) == 0 {
		return ""
	}

	type treeNode struct {
		name     string
		children []*treeNode
	}

	root := &treeNode{name: "", children: []*treeNode{}}
	for _, file := range files {
		parts := strings.Split(file, "/")
		currentNode := root
		for _, part := range parts {
			found := false
			for _, child := range currentNode.children {
				if child.name == part {
					currentNode = child
					found = true
					break
				}
			}
			if !found {
				newNode := &treeNode{name: part, children: []*treeNode{}}
				currentNode.children = append(currentNode.children, newNode)
				currentNode = newNode
			}
		}
	}

	var sb strings.Builder
	var dfs func(*treeNode, string, bool, int)
	dfs = func(n *treeNode, prefix string, isLast bool, depth int) {
		if depth == 0 {
			sb.WriteString(n.name + "\n")
		} else {
			branch := "├── "
			if isLast {
				branch = "└── "
			}
			sb.WriteString(prefix + branch + n.name + "\n")
		}

		childPrefix := prefix
		if depth >= 1 {
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│   "
			}
		}

		numChildren := len(n.children)
		for i, child := range n.children {
			isLastChild := i == numChildren-1
			dfs(child, childPrefix, isLastChild, depth+1)
		}
	}

	for i, child := range root.children {
		isLast := i == len(root.children)-1
		dfs(child, "", isLast, 0)
	}

	return sb.String()
}

func buildCommand() string {
	args := os.Args
	return strings.Join(args, " ")

}
