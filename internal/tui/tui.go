package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

const (
	extensionsInput = iota
	directoriesInput
	skipInput
	outputInput
)

var (
	focusedLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true)
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#00FF00")).
				Padding(0, 1)
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

func initialModel() model {

	dirs, err := os.ReadDir(".")
	availableDirs := []string{}
	if err == nil {
		for _, entry := range dirs {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				availableDirs = append(availableDirs, entry.Name())
			}
		}
	}

	var inputs []textinput.Model

	ti := textinput.New()
	ti.Placeholder = "py,js,go,rs,ts"
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Extensions: "
	inputs = append(inputs, ti)

	ti = textinput.New()
	ti.Placeholder = strings.Join(availableDirs, ",")
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Directories: "
	inputs = append(inputs, ti)

	ti = textinput.New()
	ti.Placeholder = "node_modules,*.tmp"
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Skip patterns: "
	inputs = append(inputs, ti)

	ti = textinput.New()
	ti.Placeholder = "llm.md"
	ti.CharLimit = 100
	ti.Width = 60
	ti.Prompt = "Output file: "
	ti.SetValue("llm.md")
	inputs = append(inputs, ti)

	return model{
		inputs:        inputs,
		focusIndex:    0,
		config:        Config{},
		availableDirs: availableDirs,
	}
}

func (m model) Init() tea.Cmd {
	return m.inputs[m.focusIndex].Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				m.processForm()
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			for i := range m.inputs {
				m.inputs[i].Blur()
			}

			cmds = append(cmds, m.inputs[m.focusIndex].Focus())
			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	input, cmd := m.inputs[m.focusIndex].Update(msg)
	m.inputs[m.focusIndex] = input
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *model) processForm() {
	if strings.TrimSpace(m.inputs[extensionsInput].Value()) == "" {
		m.err = fmt.Errorf("no extensions provided")
		return
	}

	m.err = nil

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
		m.config.OutputFile = "llm.md"
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
		Foreground(lipgloss.Color("#00FF00")).
		Padding(0, 1).
		Render("Shout - Project Dump for LLMs")

	var b strings.Builder
	b.WriteString(title + "\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("You must provide at least one extension before continuing.\n\n"))
	}

	for i, input := range m.inputs {
		prompt := input.Prompt
		content := input.View()[len(prompt):]

		if i == m.focusIndex {
			label := focusedLabelStyle.Render(prompt)
			boxed := focusedBorderStyle.Render(content)
			b.WriteString(label + "\n" + boxed + "\n\n")
		} else {
			b.WriteString(input.View() + "\n\n")
		}
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Tab/Shift+Tab: Navigate • Enter: Submit • Esc/Ctrl+C: Quit")
	b.WriteString(help)

	return b.String()
}

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
