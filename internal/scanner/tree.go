package scanner

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type treeNode struct {
	name     string
	children map[string]*treeNode
	isFile   bool
}

func (s *Scanner) generateTree(dirs []string) (string, error) {
	if len(dirs) == 0 {
		return "", errors.New("empty dirs")
	}

	root := &treeNode{
		name:     "",
		children: make(map[string]*treeNode),
		isFile:   false,
	}

	for _, dir := range dirs {
		if dir == "" {
			continue
		}

		absPath, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		err = filepath.Walk(absPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			relPath, err := filepath.Rel(absPath, path)
			if err != nil {
				return nil
			}

			norPath := strings.ReplaceAll(relPath, "\\", "/")
			components := strings.Split(norPath, "/")
			current := root

			if dir != "." {
				dirName := filepath.Base(dir)
				if _, exists := current.children[dirName]; !exists {
					current.children[dirName] = &treeNode{
						name:     dirName,
						children: make(map[string]*treeNode),
						isFile:   false,
					}
				}
				current = current.children[dirName]
			}

			for i, c := range components {
				if c == "" {
					continue
				}

				if _, exists := current.children[c]; !exists {
					isFile := i == len(components)-1
					current.children[c] = &treeNode{
						name:     c,
						children: make(map[string]*treeNode),
						isFile:   isFile,
					}
				}
				current = current.children[c]
			}
			return nil
		})
		if err != nil {
			continue
		}
	}

	var builder strings.Builder
	buildTreeString(root, &builder, "", true)
	return builder.String(), nil
}

func buildTreeString(node *treeNode, builder *strings.Builder, prefix string, isLast bool) {
	// Don't print the root node
	if node.name != "" {
		builder.WriteString(prefix)
		if isLast {
			builder.WriteString("└── ")
		} else {
			builder.WriteString("├── ")
		}
		builder.WriteString(node.name)
		builder.WriteString("\n")
	}

	// Sort children for consistent output
	var children []string
	for name := range node.children {
		children = append(children, name)
	}
	sort.Strings(children)

	// Process children
	for i, childName := range children {
		child := node.children[childName]
		childIsLast := i == len(children)-1
		newPrefix := prefix
		if node.name != "" {
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
		}
		buildTreeString(child, builder, newPrefix, childIsLast)
	}
}

func (s *Scanner) generateMetaTree(files []string) string {
	if len(files) == 0 {
		return ""
	}

	root := &treeNode{
		name:     "",
		children: make(map[string]*treeNode),
		isFile:   false,
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	for _, file := range files {

		relPath := file
		if cwd != "" {
			if rp, err := filepath.Rel(cwd, file); err == nil {
				relPath = rp
			}
		}

		nPath := strings.ReplaceAll(relPath, "\\", "/")
		components := strings.Split(nPath, "/")
		current := root

		for i, c := range components {
			if c == "" {
				continue
			}
			if _, exists := current.children[c]; !exists {
				isFile := i == len(components)-1
				current.children[c] = &treeNode{
					name:     c,
					children: make(map[string]*treeNode),
					isFile:   isFile,
				}
			}
			current = current.children[c]
		}
	}
	var builder strings.Builder
	buildTreeString(root, &builder, "", true)

	return builder.String()
}
