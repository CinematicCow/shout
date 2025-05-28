package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return m.inputs[m.focusIndex].Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeyMsg(m, msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func handleKeyMsg(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.quitting = true
		return m, tea.Quit

	case "tab", "shift+tab", "enter", "up", "down":
		return handleNavigation(m, msg.String())
	}

	index := m.focusIndex
	input, cmd := m.inputs[index].Update(msg)
	m.inputs[index] = input
	return m, cmd
}

func handleNavigation(m model, key string) (tea.Model, tea.Cmd) {
	if key == "enter" && m.focusIndex == len(m.inputs)-1 {
		m.processForm()
		return m, tea.Quit
	}

	if key == "up" || key == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex >= len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs) - 1
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) processForm() {
	if strings.TrimSpace(m.inputs[extensionsInput].Value()) == "" {
		m.err = fmt.Errorf("no extensions provided")
		return
	}

	m.err = nil
	m.config.Extensions = splitCommaString(m.inputs[extensionsInput].Value())
	m.config.Directories = getDirectories(m.inputs[directoriesInput].Value())
	m.config.SkipPatterns = splitCommaString(m.inputs[skipInput].Value())
	m.config.OutputFile = getOutputFile(m.inputs[outputInput].Value())
	m.config.Meta = getMetaOption(m.inputs[metaInput].Value())
}

func getDirectories(input string) []string {
	if val := input; val != "" {
		return splitCommaString(val)
	}
	return []string{"."}
}

func getOutputFile(input string) string {
	if val := input; val != "" {
		return val
	}
	return "llm.md"
}

func getMetaOption(input string) bool {
	metaVal := strings.ToLower(strings.TrimSpace(input))
	return metaVal == "y" || metaVal == "yes"
}
