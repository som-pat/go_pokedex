package main

import (
	"fmt"
	// "math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/som-pat/poke_dex/imagegen"
	"github.com/som-pat/poke_dex/internal/config"
	"github.com/som-pat/poke_dex/internal/replinternal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)
// var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
type btBaseModel struct {
	textInput    	textinput.Model
	output       	string
	cfgState     	*config.ConfigState
	locationList 	*Paginatedlisting
	PokemonList  	*Paginatedlisting
	UserInv			*UserInventory
	showLoc		 	bool
	showPoke	 	bool
	battlestate  	bool
	battleTarget 	string
	showmsg			bool
	selCom          int
	commands		[]string
	attributes		[]string
	selattr			int
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
		UserInv: 	  &UserInventory{},	
		showLoc: 	  false,
		showPoke: 	  false,
		battlestate:  false,
		battleTarget: "",
		showmsg:      false,
		selCom: 	  0,
		commands: 	  []string{"Attack", "Item", "Switch", "Catch", "Escape"},
		attributes:   []string{},
		selattr: 	  0,
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
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "left":
			if m.battlestate && m.selattr >0{
				m.selattr--
			}
		case "right":
			if m.battlestate && m.selattr<len(m.attributes)-1{
				m.selattr ++
			}
		case "up":
			if m.battlestate && m.selCom>0{
				m.selCom --
			}
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
			if m.battlestate && m.selCom <len(m.commands)-1{
				m.selCom++
			}
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
				if m.UserInv == nil{
					m.UserInv = InitInventoryListing()					
				}
				InventoryView(m.cfgState,m.UserInv,"fill",0)
				m.battlestate = true
				m.showLoc = false
				m.showPoke = true
				m.showmsg = true
				m.output,m.PokemonList.Items = processCommand(m.textInput.Value(),m.cfgState)
				m.battleTarget = m.PokemonList.Items[0]
				m.PokemonList.count = len(m.PokemonList.Items)
				
				m.textInput.SetValue("")
				return m, tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
					return clearMessage{}
				})
			}else if m.battlestate{

				switch m.commands[m.selCom] {
				case "Attack":
					m.output = "player attacks!"
				case "Item":
					m.attributes = m.UserInv.ItemName
					m.UserInv.ItemDescriptions=nil
					InventoryView(m.cfgState, m.UserInv,"item",m.selattr)					
					m.output= fmt.Sprintf("Switching from the %d available ones,at %d",len(m.attributes),m.selattr)
				case "Switch":					
					m.attributes = m.UserInv.PokeName
					m.UserInv.PokeDescriptions = nil
					InventoryView(m.cfgState,m.UserInv,"switch",m.selattr)
					m.output = fmt.Sprintf("Switching from the %d available ones,at %d",len(m.attributes),m.selattr)
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

func InventoryView(cfgState *config.ConfigState,uvin *UserInventory,swcase string, index int){
	if len(cfgState.PokemonCaught) == 0 && len(cfgState.ItemsHeld) == 0{
		return 
	}	
	switch swcase{
	case "fill":
		for _,poke := range cfgState.PokemonCaught{
			uvin.PokeName = append(uvin.PokeName, poke.Name)
		}
		for _,item := range cfgState.ItemsHeld{
			uvin.ItemName = append(uvin.ItemName, item.Name)
		}
	case "switch":
		for _, poke := range cfgState.PokemonCaught{
			if poke.Name == uvin.PokeName[index]{
				ascii_img, _ := imagegen.AsciiGen(poke.Sprites.FrontDefault,56)
				uvin.PokeSprite =  ascii_img
				for _,stats := range poke.Stats{
					uvin.PokeDescriptions = append(uvin.PokeDescriptions, strconv.Itoa(stats.BaseStat))
				}
			}
		}	
	case "item":
		for _,item := range cfgState.ItemsHeld{
			uvin.ItemDescriptions= append(uvin.ItemDescriptions, item.Category.Name)
			uvin.ItemDescriptions = append(uvin.ItemDescriptions, item.EffectEntries[0].ShortEffect)
		}

	}	
}

type clearMessage struct{}

type UserInventory struct{
	PokeSprite				string
	PokeName				[]string
	ItemName				[]string
	MoveName				[]string
	PokeDescriptions		[]string
	ItemDescriptions		[]string
	MoveDescriptions		[]string

}

func InitInventoryListing() *UserInventory {
	return &UserInventory{
		PokeSprite:  	  "",
		PokeName: 		  []string{},
		ItemName: 		  []string{},
		MoveName: 		  []string{},	
		PokeDescriptions: []string{},
		ItemDescriptions: []string{},
		MoveDescriptions: []string{},
	}
}



func HealthBar(curHP,maxHP int)string{
	hpercent := float64(curHP)/float64(maxHP)
	barwidth := 10

	filledWidth :=int(hpercent*float64(barwidth))
	emptyWidth:= barwidth - filledWidth
	
	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#004400"))
    emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Background(lipgloss.Color("#222222"))

	filled := filledStyle.Render(strings.Repeat("█", filledWidth))
	empty :=  emptyStyle.Render(strings.Repeat(" ",emptyWidth))
	hpercent *= 100
	return fmt.Sprintf("%s%s %d%%", filled, empty, int(hpercent))

}

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
		attributes := make([]string, len(m.attributes))
		for i,attr := range m.attributes{
			style := lipgloss.NewStyle().Bold(true)
			if i == m.selattr{
				style = style.Foreground(lipgloss.Color("220")).Underline(true)
			}else{
				style = style.Foreground(lipgloss.Color("240"))
			}
			attributes[i] = style.Render("[" + attr+ "]")
		}
		var plstats string
		var tlstats string

		commandBoxContent := lipgloss.JoinVertical(lipgloss.Top, commands...)
		commandAttrContent := lipgloss.JoinHorizontal(lipgloss.Top, attributes...)

		commandListBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 0).Width(28).Align(lipgloss.Center).Render(commandBoxContent)
		commandAttributes := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0,0).Width(100).Align(lipgloss.Center).Render(commandAttrContent)

		commandBox := lipgloss.JoinHorizontal(lipgloss.Left,commandListBox,commandAttributes)

		enemyViewBox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0,0,0,0).Width(64).Height(15)
		enemyStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0,0,0,0).Width(64).Height(3)
		playerViewBox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0, 0, 0).Width(64).Height(15)
		playerStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0,0,0,0).Width(64).Height(3)
		

		
		if len(m.PokemonList.Items)>4 {
			hp,err := strconv.Atoi(m.PokemonList.Items[3])
			if err!=nil{return "Error Converting"}
			enemyhb := HealthBar(hp,hp)
			lvlhb := lipgloss.JoinHorizontal(lipgloss.Left,m.PokemonList.Items[2]," ",enemyhb)
			tlstats = lipgloss.JoinVertical(lipgloss.Top,cases.Title(language.Und, cases.NoLower).String((m.PokemonList.Items[0])),lvlhb)
		}else{
			tlstats = lipgloss.JoinVertical(lipgloss.Top,"Enemy Pokemon")
		}		

		if len(m.UserInv.PokeDescriptions) > 0 {
			hp,err := strconv.Atoi(m.UserInv.PokeDescriptions[0])
			if err!=nil{return "Error Converting"}
			playerhb := HealthBar(hp, hp)
			plstats = lipgloss.JoinVertical(lipgloss.Top, cases.Title(language.Und, cases.NoLower).String((m.UserInv.PokeName[m.selattr])), playerhb)
		} else {
			plstats = lipgloss.JoinVertical(lipgloss.Top, "Player 1")
		}

		topLeftPokemon := lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(m.PokemonList.Items[1])
		topLeftStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(tlstats)
		
		bottomRightPokemon := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CFFF")).Render(m.UserInv.PokeSprite)
		bottomRightStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CFFF")).Render(plstats)

		
		// Display Pokémon and commands
		enemyPokemonView := lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.PlaceHorizontal(0,lipgloss.Left,topLeftPokemon))
		
		playerPokemonView := lipgloss.JoinVertical(lipgloss.Bottom,
			lipgloss.PlaceHorizontal(64,lipgloss.Right,bottomRightPokemon))
		
		enemyBox := lipgloss.JoinVertical(lipgloss.Top,enemyViewBox.Render(enemyPokemonView),
					playerStatusBox.Render(topLeftStatus))
		playerBox := lipgloss.JoinVertical(lipgloss.Right, enemyStatusBox.Render(bottomRightStatus),
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
		"[Press 'esc' to quit, 'up'/'down'/'right'/'left' to navigate, 'enter' to confirm]",
	)
}


func processCommand(input string, cfgState *config.ConfigState) (string, []string) {
	input = strings.TrimSpace(input)
	fmt.Printf("Input given %s ", input)
	res, lis := replinternal.ReplInput(cfgState, input)
	return res, lis
}
