package storymode

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/som-pat/poke_dex/app"
	"github.com/som-pat/poke_dex/internal/config"
	"github.com/som-pat/poke_dex/internal/replinternal"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type SMStoryModel struct{
	textInput    	textinput.Model
	locationList 	*Paginatedlisting
	PokemonList  	*Paginatedlisting
	output       	string
	cfgstate 		*config.ConfigState
	Navigator 	  	*app.AppNavigator
	width     		int
	height	  		int
	showLoc		 	bool
	showPoke	 	bool
	showinspect		bool
}

type Paginatedlisting struct {
	Items        []string
	count        int
	selectedIndex int
}


func StoryInput(cfgState *config.ConfigState,  navigator *app.AppNavigator) SMStoryModel{
	cti := textinput.New()
	cti.Placeholder = "[Press 'esc' to quit, 'up'/'down'/'right'/'left' to navigate, 'enter' to confirm]"
	cti.Focus()
	cti.CharLimit = 200
	cti.Width = 200
	return SMStoryModel{
		textInput:    	cti,
		output:       	"Welcome to PokeCLI!\nType 'explore' to search for Pokémon, 'catch' to catch one, or 'quit' to exit.",
		cfgstate:		cfgState,
		Navigator: 		navigator,
		locationList: 	InitPaginatedListing(20),
		PokemonList:  	InitPaginatedListing(40),
		height: 		0,
		width:			0,
		showLoc: 		false,
		showPoke: 		false,
	}
}

func InitPaginatedListing(capmap int) *Paginatedlisting {
	return &Paginatedlisting{
		Items:        make([]string, 0, capmap),
		count:        0,
		selectedIndex: 0,
	}
}

func (m SMStoryModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SMStoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg:= msg.(type){
	case tea.WindowSizeMsg:
		m.width  = msg.Width
		m.height = msg.Height
	
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.Navigator.GoToMenu(), nil
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
				if m.PokemonList.count  == 0 {
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
			}else if strings.HasPrefix(m.textInput.Value(), "catch") {
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
				m.output, m.PokemonList.Items = processCommand(m.textInput.Value(), m.cfgstate)
				m.PokemonList.count = len(m.PokemonList.Items)
				m.showinspect = false
				m.showPoke =true
				m.showLoc = false
				// m.battlestate = false
				// m.showmsg = false
				// m.attackstate = false
				
				m.textInput.SetValue("") 
			}else {
				m.output, m.locationList.Items = processCommand(m.textInput.Value(), m.cfgstate)
				m.locationList.count = len(m.locationList.Items)
				m.showinspect =false
				part := strings.Split(m.textInput.Value(), " ")
				if part[0] =="inspect"{m.showinspect = true}
				m.showLoc = true
				m.showPoke =false				
				// m.battlestate = false
				// m.showmsg = false
				// m.attackstate = false
				m.textInput.SetValue("")
			}
		}
	m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m SMStoryModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("35"))
	screenHeight := 44	
	outputLines := strings.Split(m.output, "\n")
	outputHeight := len(outputLines)	
	paddingLines := screenHeight - outputHeight - 5 
	
	padding := strings.Repeat("\n", paddingLines)
	
	return fmt.Sprintf(
		"%s\n%s%s %s",
		headerStyle.Render(m.output),
		padding,
		"Pokedex", m.textInput.View(),
	)
	// return "Story Mode: " + m.cfgstate.PlayerName + "\nPress ESC to return to the menu."
}

func processCommand(input string, cfgState *config.ConfigState) (string, []string) {
	input = strings.TrimSpace(input)
	fmt.Printf("Input given %s ", input)
	res, lis := replinternal.ReplInput(cfgState, input)
	return res, lis
}