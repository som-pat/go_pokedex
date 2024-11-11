package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type btBaseModel struct {
	textInput    textinput.Model
	output       string
	cfgState     *ConfigState
	locationList *Paginatedlisting
	PokemonList  *Paginatedlisting
	showLoc		 bool
	showPoke	 bool
}

type Paginatedlisting struct {
	Items        []string
	count        int
	selectedIndex int
}


func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: 
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func takeInput(cfgState *ConfigState) btBaseModel {
	ti := textinput.New()
	// ti.Placeholder = "pokedex > "
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	return btBaseModel{
		textInput:    ti,
		output:       "Welcome to PokeCLI!\nType 'explore' to search for Pokémon, 'catch' to catch one, or 'quit' to exit.",
		cfgState:     cfgState,
		locationList: InitPaginatedListing(20),
		PokemonList:  InitPaginatedListing(40),
		showLoc: 	  false,
		showPoke: 	  false,
	}
}

func InitPaginatedListing(capmap int) *Paginatedlisting {
	return &Paginatedlisting{
		Items:        make([]string, 0, capmap),
		count:        0,
		selectedIndex: 0,
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
		case "up":
			if strings.HasPrefix(m.textInput.Value(), "explore") || strings.HasPrefix(m.textInput.Value(),"inspect"){
				if m.locationList.count == 0 {
					m.output = "Execute map first"
					m.textInput.SetValue("")
				}else{
					if m.locationList.selectedIndex>0{
						m.locationList.selectedIndex --
					}else{
						m.locationList.selectedIndex = m.locationList.count -1
					}
				}
				fmt.Print(m.locationList.selectedIndex,m.showLoc)
				if m.showLoc && strings.HasPrefix(m.textInput.Value(), "explore") {
					selLoc := m.locationList.Items[m.locationList.selectedIndex]
					m.textInput.SetValue(fmt.Sprintf("explore %s",selLoc))
				}else if m.showLoc && strings.HasPrefix(m.textInput.Value(),"inspect"){
					selLoc := m.locationList.Items[m.locationList.selectedIndex]
					m.textInput.SetValue(fmt.Sprintf("inspect %s",selLoc))
				}
			} else if strings.HasPrefix(m.textInput.Value(), "catch") {
				if m.PokemonList.count == 0 {
					m.output = "Selct a region to explore first"
					m.textInput.SetValue("")
				}else{
					if m.PokemonList.selectedIndex > 0{
						m.PokemonList.selectedIndex --
					}else{
						m.PokemonList.selectedIndex = m.PokemonList.count -1
					}
				}
				fmt.Printf("Not gonna catch'em all,%d",m.PokemonList.selectedIndex)
				if m.showPoke{
					selPok := m.PokemonList.Items[m.PokemonList.selectedIndex]
					m.textInput.SetValue(fmt.Sprintf("catch %s",selPok))
				}
			}
		case "down":
			if strings.HasPrefix(m.textInput.Value(), "explore") || strings.HasPrefix(m.textInput.Value(),"inspect"){
				if m.locationList.count == 0 {
					m.output = "Execute map first"
					m.textInput.SetValue("")
				}else{
					if m.locationList.selectedIndex < m.locationList.count -1{
						m.locationList.selectedIndex ++
					}else{
						m.locationList.selectedIndex = 0
					}
				}
				fmt.Print(m.locationList.selectedIndex,m.showLoc)
				if m.showLoc && strings.HasPrefix(m.textInput.Value(), "explore") {
					selLoc := m.locationList.Items[m.locationList.selectedIndex]
					m.textInput.SetValue(fmt.Sprintf("explore %s",selLoc))
				}else if m.showLoc && strings.HasPrefix(m.textInput.Value(),"inspect"){
					selLoc := m.locationList.Items[m.locationList.selectedIndex]
					m.textInput.SetValue(fmt.Sprintf("inspect %s",selLoc))
				}
			} else if strings.HasPrefix(m.textInput.Value(), "catch") {
				if m.PokemonList.count == 0 {
					m.output = "Selct a region to explore first"
					m.textInput.SetValue("")
				} else{
					if m.PokemonList.selectedIndex < m.PokemonList.count -1{
						m.PokemonList.selectedIndex ++
					}else{
						m.PokemonList.selectedIndex = 0 
					}
				}

				fmt.Printf("Not gonna catch'em all,%d",m.PokemonList.selectedIndex)
				if m.showPoke{
				selPok := m.PokemonList.Items[m.PokemonList.selectedIndex]
				m.textInput.SetValue(fmt.Sprintf("catch %s",selPok))
				}
			} 
		case "enter":
			m.locationList.selectedIndex = 0
			m.PokemonList.selectedIndex = 0
			m.output =""
			if strings.HasPrefix(m.textInput.Value(), "explore") || strings.HasPrefix(m.textInput.Value(), "scout") {				
				m.output, m.PokemonList.Items = processCommand(m.textInput.Value(), m.cfgState)
				m.PokemonList.count = len(m.PokemonList.Items)
				m.showPoke =true
				m.showLoc = false
				m.textInput.SetValue("") 
			}else {
				m.output, m.locationList.Items = processCommand(m.textInput.Value(), m.cfgState)
				m.locationList.count = len(m.locationList.Items)
				m.showPoke =false
				m.showLoc = true
				m.textInput.SetValue("")
			}
		}

	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m btBaseModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("35"))
	screenHeight := 41 // Adjust this based on your terminal height if necessary

	// Count lines in the output
	outputLines := strings.Split(m.output, "\n")
	outputHeight := len(outputLines)

	// Calculate padding needed to push the prompt to the bottom
	paddingLines := screenHeight - outputHeight - 5 // Extra space for the prompt and instructions

	padding := strings.Repeat("\n", paddingLines)

	return fmt.Sprintf(
		"%s\n%s\n%s %s\n\n%s\n",
		headerStyle.Render(m.output),
		padding,
		"Pokedex", m.textInput.View(),
		"[Press 'q' to quit, 'up'/'down' to navigate, 'enter' to confirm]",
	)
}


func processCommand(input string, cfgState *ConfigState) (string, []string) {
	input = strings.TrimSpace(input)
	fmt.Printf("Input given %s ", input)
	res, lis := repl_input(cfgState, input)
	return res, lis
}
