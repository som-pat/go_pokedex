package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/som-pat/poke_dex/internal/config"
	"github.com/som-pat/poke_dex/internal/replinternal"
	"github.com/som-pat/poke_dex/imagegen"
)

type btBaseModel struct {
	textInput    	textinput.Model
	output       	string
	cfgState     	*config.ConfigState
	locationList 	*Paginatedlisting
	PokemonList  	*Paginatedlisting
	showLoc		 	bool
	showPoke	 	bool
	battlestate  	bool
	battleTarget 	string
	showmsg			bool
	selCom          int
	commands		[]string
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

func takeInput(cfgState *config.ConfigState) btBaseModel {
	ti := textinput.New()
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
		battlestate:  false,
		battleTarget: "",
		showmsg:      false,
		selCom: 	  0,
		commands: 	  []string{"Attack", "Item", "Switch", "Catch", "Escape"},
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
		case "left":
			if m.battlestate && m.selCom >0{
				m.selCom--
			}
		case "right":
			if m.battlestate && m.selCom<len(m.commands)-1{
				m.selCom ++
			}
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
				m.battlestate = false
				m.showmsg = false
				m.textInput.SetValue("") 
			}else if strings.HasPrefix(m.textInput.Value(),"battle") {
				m.battlestate = true
				m.showLoc = false
				m.showPoke = true
				m.showmsg = true
				m.output,m.PokemonList.Items = processCommand(m.textInput.Value(),m.cfgState)
				m.battleTarget = m.PokemonList.Items[0]
				m.PokemonList.count = len(m.PokemonList.Items)
				m.textInput.SetValue("")
				return m, tea.Tick(5*time.Second, func(_ time.Time) tea.Msg {
					// Clear the initial message after 1 second
					return clearMessage{}
				})
			}else if m.battlestate{

				switch m.commands[m.selCom] {
				case "Attack":
					m.output = "player attacks!"
				case "Item":					
					m.output, _ = InventoryView(m.cfgState,"item")
				case "Switch":					
					m.output,m.PokemonList.Items = InventoryView(m.cfgState, "switch")
				case "Catch":
					m.output = "Catching Pokemon, choose the balls!"
				case "Escape":
					m.output = "You escaped the battle!"
					m.battlestate = false
				}

				m.textInput.SetValue("")			
			
				
			}else {
				m.output, m.locationList.Items = processCommand(m.textInput.Value(), m.cfgState)
				m.locationList.count = len(m.locationList.Items)
				m.showPoke =false
				m.showLoc = true
				m.battlestate = false
				m.showmsg = false
				m.textInput.SetValue("")
			}		
		}
	case clearMessage:
		// Clear the initial battle message and switch to the battle HUD
		m.showmsg = false
		m.output = ""
		return m, cmd
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func InventoryView(cfgState *config.ConfigState,swcase string ) (string,[]string){
	var inventory strings.Builder
	var contianer []string
	if swcase == "switch"{
		if len(cfgState.PokemonCaught)==0{
			inventory.WriteString("Try with balls caus you dont have any poke to switch with")
			return inventory.String(),nil
		}
		for _, poke := range cfgState.PokemonCaught{
			inventory.WriteString(fmt.Sprintf("%s \n",poke.Name))
			ascii_img, _ := imagegen.AsciiGen(poke.Sprites.FrontDefault,16)
			contianer = append(contianer, ascii_img)
			}
	}else{
		if len(cfgState.ItemsHeld)==0{
			inventory.WriteString("No balls to try")
			return inventory.String(),nil
		}
		for _, item := range cfgState.ItemsHeld{
			inventory.WriteString(fmt.Sprintf("%s \n",item.Name))
		}
	}		
	return inventory.String(),contianer
}

type clearMessage struct{}

func (m btBaseModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("35"))
	screenHeight := 41 

	
	outputLines := strings.Split(m.output, "\n")
	outputHeight := len(outputLines)

	
	paddingLines := screenHeight - outputHeight - 5 

	padding := strings.Repeat("\n", paddingLines)


	if m.battlestate{
		if m.showmsg{
			return fmt.Sprintf("%s\n\n%s",m.output, m.textInput.View())
		}

		// Battle HUD layout
		commands := make([]string, len(m.commands))
		for i, cmd := range m.commands {
			style := lipgloss.NewStyle().Bold(true)
			if i == m.selCom {
				style = style.Foreground(lipgloss.Color("220")).Underline(true)
			} else {
				style = style.Foreground(lipgloss.Color("240"))
			}
			commands[i] = style.Render("[" + cmd + "]")
		}

		commandBoxContent := lipgloss.JoinHorizontal(lipgloss.Top, commands...)
		commandBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 2).Width(128).Align(lipgloss.Center).Render(commandBoxContent)


		enemyViewBox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0,0,0,0).Width(64).Height(15)
		enemyStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0,0,0,0).Width(64).Height(3)
		playerViewBox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0, 0, 0).Width(64).Height(15)
		playerStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0,0,0,0).Width(64).Height(3)
		

		topLeftPokemon := lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(m.PokemonList.Items[1])
		topLeftStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(m.PokemonList.Items[0])
		bottomRightPokemon := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CFFF")).Render(m.PokemonList.Items[1])
		bottomRightStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CFFF")).Render(m.PokemonList.Items[0])

		
		// Display Pokémon and commands
		enemyPokemonView := lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.PlaceHorizontal(0,lipgloss.Left,topLeftPokemon))
		
		playerPokemonView := lipgloss.JoinVertical(lipgloss.Bottom,
			lipgloss.PlaceHorizontal(64,lipgloss.Right,bottomRightPokemon))
		
		enemyBox := lipgloss.JoinVertical(lipgloss.Top,enemyViewBox.Render(enemyPokemonView),
					playerStatusBox.Render(topLeftStatus))
		playerBox := lipgloss.JoinVertical(lipgloss.Top, enemyStatusBox.Render(bottomRightStatus),
					 playerViewBox.Render(playerPokemonView))			
		
		
		battleBox :=lipgloss.JoinHorizontal(lipgloss.Left,enemyBox,
					playerBox)
		// pokemonView := lipgloss.JoinVertical(lipgloss.Left,
		// 	lipgloss.PlaceHorizontal(64, lipgloss.Left, topLeftPokemon),
		// 	lipgloss.NewStyle().Height(0).Render(""), // Spacer for middle
		// 	lipgloss.PlaceHorizontal(128, lipgloss.Right, bottomRightPokemon),
		// )

		// commands := lipgloss.JoinHorizontal(lipgloss.Top,
		// 	lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render("[Attack]"),
		// 	lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")).Render("[Item]"),
		// 	lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).Render("[Switch]"),
		// 	lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Render("[Catch]"),
		// 	lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render("[Escape]"),
		// )

		return lipgloss.JoinVertical(lipgloss.Left,
			battleBox,
			commandBox,
			m.output,
			"\n", m.textInput.View(),
		)
	}
	

	return fmt.Sprintf(
		"%s\n%s\n%s %s\n\n%s\n",
		headerStyle.Render(m.output),
		padding,
		"Pokedex", m.textInput.View(),
		"[Press 'q' to quit, 'up'/'down' to navigate, 'enter' to confirm]",
	)
}


func processCommand(input string, cfgState *config.ConfigState) (string, []string) {
	input = strings.TrimSpace(input)
	fmt.Printf("Input given %s ", input)
	res, lis := replinternal.ReplInput(cfgState, input)
	return res, lis
}
