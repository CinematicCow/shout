package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

const (
	extensionsInput = iota
	directoriesInput
	skipInput
	outputInput
	metaInput
)

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

func initialModel() model {
	availableDirs := getAvailableDirs()
	return model{
		inputs:        createInputs(availableDirs),
		focusIndex:    0,
		config:        Config{},
		availableDirs: availableDirs,
	}
}

func getAvailableDirs() []string {
	dirs, err := os.ReadDir(".")
	if err != nil {
		return []string{}
	}

	availableDirs := []string{}
	for _, entry := range dirs {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			availableDirs = append(availableDirs, entry.Name())
		}
	}
	return availableDirs
}

func createInputs(availableDirs []string) []textinput.Model {
	inputs := make([]textinput.Model, 5)

	inputs[extensionsInput] = createTextInput("Extensions: ", "py,js,go,rs,ts", 100, 60)
	inputs[directoriesInput] = createTextInput("Directories: ", strings.Join(availableDirs, ","), 100, 60)
	inputs[skipInput] = createTextInput("Skip patterns: ", "node_modules,*.tmp", 100, 60)
	
	outInput := createTextInput("Output file: ", "llm.md", 100, 60)
	outInput.SetValue("llm.md")
	inputs[outputInput] = outInput

	metaIn := createTextInput("Generate Meta file? ", "Y/n", 1, 60)
	inputs[metaInput] = metaIn

	return inputs
}

func createTextInput(prompt, placeholder string, charLimit, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	ti.Width = width
	ti.Prompt = prompt
	return ti
}
