package scanner

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func (s *Scanner) generateTree(dirs []string) (string, error) {
	if output, err := exec.Command("tree", append([]string{"--noreport"}, dirs...)...).Output(); err == nil {
		return string(output), nil
	}

	var sb strings.Builder
	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil || !info.IsDir() {
				return err
			}

			relPath, _ := filepath.Rel(dir, path)
			if relPath == "." {
				sb.WriteString(filepath.Base(dir) + "\n")
			} else {
				depth := strings.Count(relPath, string(filepath.Separator))
				sb.WriteString(strings.Repeat("|   ", depth) + "|-- " + filepath.Base(path) + "\n")
			}
			return nil
		})
	}
	return sb.String(), nil
}

type treeNode struct {
	name     string
	children []*treeNode
}

func (s *Scanner) generateMetaTree(files []string) string {
	if len(files) == 0 {
		return ""
	}

	root := &treeNode{}
	for _, file := range files {
		parts := strings.Split(file, "/")
		current := root
		for _, part := range parts {
			found := false
			for _, child := range current.children {
				if child.name == part {
					current = child
					found = true
					break
				}
			}
			if !found {
				newNode := &treeNode{name: part}
				current.children = append(current.children, newNode)
				current = newNode
			}
		}
	}

	var sb strings.Builder
	var dfs func(*treeNode, string, bool)
	dfs = func(n *treeNode, prefix string, isLast bool) {
		if n.name != "" {
			branch := "├── "
			if isLast {
				branch = "└── "
			}
			sb.WriteString(prefix + branch + n.name + "\n")
			if isLast {
				prefix += "    "
			} else {
				prefix += "│   "
			}
		}

		sort.Slice(n.children, func(i, j int) bool {
			return n.children[i].name < n.children[j].name
		})

		for i, child := range n.children {
			dfs(child, prefix, i == len(n.children)-1)
		}
	}

	dfs(root, "", true)
	return sb.String()
}
