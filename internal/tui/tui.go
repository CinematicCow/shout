package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Config holds the TUI configuration
type Config struct {
	Extensions   []string
	Directories  []string
	SkipPatterns []string
	OutputFile   string
}

type model struct {
	inputs        []textinput.Model
	focusIndex    int
	config        Config
	err           error
	width         int
	height        int
	quitting      bool
	availableDirs []string
}

// Input field indices
const (
	extensionsInput = iota
	directoriesInput
	skipInput
	outputInput
)

func initialModel() model {
	var inputs []textinput.Model

	// Extensions input
	ti := textinput.New()
	ti.Placeholder = "py,js,go,rs,ts"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Extensions: "
	inputs = append(inputs, ti)

	// Directories input
	ti = textinput.New()
	ti.Placeholder = "."
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Directories: "
	inputs = append(inputs, ti)

	// Skip patterns input
	ti = textinput.New()
	ti.Placeholder = "node_modules,*.tmp"
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Skip patterns: "
	inputs = append(inputs, ti)

	// Output file input
	ti = textinput.New()
	ti.Placeholder = "shout.md"
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Output file: "
	ti.SetValue("shout.md")
	inputs = append(inputs, ti)

	// List available directories
	dirs, err := os.ReadDir(".")
	availableDirs := []string{}
	if err == nil {
		for _, entry := range dirs {
			if entry.IsDir() {
				availableDirs = append(availableDirs, entry.Name())
			}
		}
	}

	return model{
		inputs:        inputs,
		focusIndex:    0,
		config:        Config{},
		availableDirs: availableDirs,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			// Cycle between inputs
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				// Process the form when Enter is pressed on the last input
				m.processForm()
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Handle character input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only update the focused input
	_, cmds[m.focusIndex] = m.inputs[m.focusIndex].Update(msg)

	return tea.Batch(cmds...)
}

func (m *model) processForm() {
	// Convert input values to config
	if val := m.inputs[extensionsInput].Value(); val != "" {
		m.config.Extensions = splitCommaString(val)
	}

	if val := m.inputs[directoriesInput].Value(); val != "" {
		m.config.Directories = splitCommaString(val)
	} else {
		m.config.Directories = []string{"."}
	}

	if val := m.inputs[skipInput].Value(); val != "" {
		m.config.SkipPatterns = splitCommaString(val)
	}

	if val := m.inputs[outputInput].Value(); val != "" {
		m.config.OutputFile = val
	} else {
		m.config.OutputFile = "shout.md"
	}
}

func splitCommaString(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Shout - Project Dump for LLMs")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Tab/Shift+Tab: Navigate • Enter: Submit • Esc/Ctrl+C: Quit")

	var b strings.Builder
	b.WriteString(title + "\n\n")

	for i, input := range m.inputs {
		b.WriteString(input.View() + "\n")

		// Show available directories under the directories input
		if i == directoriesInput && len(m.availableDirs) > 0 {
			availableDirsText := fmt.Sprintf("Available: %s", strings.Join(m.availableDirs, ", "))
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Italic(true).
				Render(availableDirsText) + "\n")
		}

		b.WriteString("\n")
	}

	b.WriteString("\n" + help)

	return b.String()
}

// StartTUI launches the TUI and returns the configuration
func StartTUI() (Config, error) {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return Config{}, err
	}

	finalM, _ := finalModel.(model)
	if finalM.quitting {
		os.Exit(0)
	}

	return finalM.config, finalM.err
}
