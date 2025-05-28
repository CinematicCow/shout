package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func StartTUI() (Config, error) {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return Config{}, err
	}

	if finalM, ok := finalModel.(model); ok {
		if finalM.quitting {
			os.Exit(0)
		}
		return finalM.config, finalM.err
	}

	return Config{}, fmt.Errorf("invalid model type")
}
