package main

import (
	"fmt"
	"strings"
    "os"
	"os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

)

type btBaseModel struct{
    textInput textinput.Model
    output    string
	cfgState  *ConfigState
}

func clearScreen() {
	cmd := exec.Command("clear") // Use "cls" for Windows
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func takeInput(cfgState *ConfigState) btBaseModel{
	ti := textinput.New()
	// ti.Placeholder = "pokedex > "
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	return btBaseModel{
		textInput: ti,
		output:    "Welcome to PokeCLI!\nType 'explore' to search for Pokémon, 'catch' to catch one, or 'quit' to exit.",
		cfgState: cfgState,
	}
}

func (m btBaseModel) Init() tea.Cmd {
    clearScreen()
	return textinput.Blink
}

func (m btBaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.output = processCommand(m.textInput.Value(),m.cfgState)
			m.textInput.SetValue("")
		}
	
	}

	 
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m btBaseModel) View() string {
	return fmt.Sprintf(
		"%s\n\nPokedex%s\n\n%s\n",
		m.output,
		m.textInput.View(),
		"[Press 'q' to quit]",
	)
}

func processCommand(input string,cfgState *ConfigState) string {
	input = strings.TrimSpace(input)	
	fmt.Printf("Input given %s ", input )
	res := repl_input(cfgState, input)
	return res


	// switch input {
	// case "explore":
	// 	return "Exploring... You find a wild Pokémon!"
	// case "catch":
	// 	return "Attempting to catch Pokémon... Success!"
	// case "help":
	// 	return "Available commands: explore, catch, help, quit"
	// case "quit":
	// 	return "Goodbye!"
	// default:
	// 	return "Unknown command. Type 'help' for available commands."
	// }
}

