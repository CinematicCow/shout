package scanner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Git struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
	Files   []string
}

func GetGit(limit int) ([]Git, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git command not found")
	}

	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		return nil, fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%ad|%s", "--date=iso", "-n", fmt.Sprintf("%d", limit))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git history: %w", err)
	}

	commits := []Git{}
	lines := strings.SplitSeq(string(output), "\n")

	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}
		hash := parts[0]
		author := parts[1]
		dateStr := parts[2]
		message := parts[3]

		date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			continue
		}

		filesCmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", hash)
		filesOutput, err := filesCmd.Output()
		if err != nil {
			continue
		}

		files := strings.Split(strings.TrimSpace(string(filesOutput)), "\n")
		if len(files) == 1 && files[0] == "" {
			files = []string{}
		}

		commits = append(commits, Git{
			Hash:    hash,
			Author:  author,
			Date:    date,
			Message: message,
			Files:   files,
		})
	}

	return commits, nil
}

func FormatGit(commits []Git) string {
	if len(commits) == 0 {
		return "No git history \n"
	}

	var buf bytes.Buffer
	buf.WriteString("## Version History\n\n")

	for _, commit := range commits {
		buf.WriteString(fmt.Sprintf("### %s (%s)\n", commit.Hash[:7], commit.Date.Format("2006-01-02")))
		buf.WriteString(fmt.Sprintf("**Author:** %s\n", commit.Author))
		buf.WriteString(fmt.Sprintf("**Message:** %s\n", commit.Message))

		if len(commit.Files) > 0 {
			buf.WriteString("**Files changed:**\n")
			for _, file := range commit.Files {
				buf.WriteString(fmt.Sprintf("- %s\n", file))
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
