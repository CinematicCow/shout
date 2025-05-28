package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(renderTitle())

	if m.err != nil {
		b.WriteString(renderError())
	}

	for i := range m.inputs {
		b.WriteString(renderInput(m, i))
	}

	b.WriteString(renderHelp())
	return b.String()
}

func renderTitle() string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")).
		Padding(0, 1).
		Render("Shout - Project Dump for LLMs") + "\n\n"
}

func renderError() string {
	return errorStyle.Render("You must provide at least one extension before continuing.\n\n")
}

func renderInput(m model, i int) string {
	input := m.inputs[i]
	prompt := input.Prompt
	content := input.View()[len(prompt):]

	if i == m.focusIndex {
		return renderFocusedInput(prompt, content)
	}
	return input.View() + "\n\n"
}

func renderFocusedInput(prompt, content string) string {
	label := focusedLabelStyle.Render(prompt)
	boxed := focusedBorderStyle.Render(content)
	return label + "\n" + boxed + "\n\n"
}

func renderHelp() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Tab/Shift+Tab: Navigate • Enter: Submit • Esc/Ctrl+C: Quit")
}
