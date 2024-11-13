package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"unicode"

	// "os"
	// "os/exec"
	// "runtime"
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
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
const cmdview = 5

type btBaseModel struct {
	textInput    	textinput.Model
	output       	string
	cfgState     	*config.ConfigState
	locationList 	*Paginatedlisting
	PokemonList  	*Paginatedlisting
	UserInv			*UserInventory
	batlog 			*BattleLog
	MoAt			*MoveAttack
	enemy			*EnemyPokemonStats
	player			*PlayerPokemonStats
	showLoc		 	bool
	showPoke	 	bool
	battlestate  	bool
	battleTarget 	string
	showmsg			bool
	selCom          int
	commands		[]string
	attributes		[]string
	selswitch		int
	selmove			int
	selitem			int
	trackcomind		int
	attackstate		bool
	turnComplete	chan bool
	moveMutex		sync.Mutex

}

type Paginatedlisting struct {
	Items        []string
	count        int
	selectedIndex int
}



func takeInput(cfgState *config.ConfigState) *btBaseModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	return &btBaseModel{
		textInput:    ti,
		output:       "Welcome to PokeCLI!\nType 'explore' to search for Pokémon, 'catch' to catch one, or 'quit' to exit.",
		cfgState:     cfgState,
		locationList: InitPaginatedListing(20),
		PokemonList:  InitPaginatedListing(40),
		UserInv: 	  &UserInventory{},
		batlog: 	  InitBattlelogIniate(),
		enemy: 		  &EnemyPokemonStats{},
		player:		  &PlayerPokemonStats{},		
		showLoc: 	  false,
		showPoke: 	  false,
		battlestate:  false,
		battleTarget: "",
		showmsg:      false,
		selCom: 	  0,
		commands: 	  []string{"Attack", "Item", "Switch", "Catch", "ViewStats", "Escape"},
		attributes:   []string{},
		selswitch: 	  0,
		selmove: 	  0,
		selitem: 	  0,
		trackcomind:  0,
		attackstate:  false,
		MoAt: 		  &MoveAttack{},
		moveMutex:	  sync.Mutex{},	
		turnComplete: make(chan bool),
	}
}

func InitPaginatedListing(capmap int) *Paginatedlisting {
	return &Paginatedlisting{
		Items:        make([]string, 0, capmap),
		count:        0,
		selectedIndex: 0,
	}
}

func (m *btBaseModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *btBaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	select {
	case <-m.turnComplete:
		m.attackstate = false
		m.output += "\nTurn completed. Choose your next move."
	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				return m, tea.Quit
			case "left":
				switch m.commands[m.selCom] {
				case "Attack":
					if m.selmove > 0 {
						m.selmove--
					}
				case "Item":
					if m.selitem > 0 {
						m.selitem--
					}
				case "Switch":
					if m.selswitch > 0 {
						m.selswitch--
					}
				}
			case "right":
				switch m.commands[m.selCom] {
				case "Attack":
					if m.selmove < len(m.UserInv.MoveName)-1 {
						m.selmove++
					}
				case "Item":
					if m.selitem < len(m.UserInv.ItemName)-1 {
						m.selitem++
					}
				case "Switch":
					if m.selswitch < len(m.UserInv.PokeName)-1 {
						m.selswitch++
					}
				}
			
			case "up":
				if m.battlestate && m.selCom>0{				
					m.selCom --
					if m.selCom < m.trackcomind{
						m.trackcomind --
					}
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
					if m.selCom >= m.trackcomind+cmdview{
						m.trackcomind++
					}
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
					m.showPoke =true
					m.showLoc = false
					m.battlestate = false
					m.showmsg = false
					m.attackstate = false
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
					m.attackstate = false
					m.output,m.PokemonList.Items = processCommand(m.textInput.Value(),m.cfgState)
					if m.enemy != nil{m.enemy = nil}
					m.enemy = InitEnpokeStats(m.PokemonList.Items)
					m.battleTarget = m.enemy.Name
					m.PokemonList.Items = nil			
					m.textInput.SetValue("")
					return m, tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
						return clearMessage{}
					})

				}else if m.battlestate && !m.attackstate{				
					switch m.commands[m.selCom] {
					case "Attack":					
						m.UserInv.MoveName =nil
						InventoryView(m.cfgState,m.UserInv,"move",m.selswitch)
						m.attributes = m.UserInv.MoveName
						m.attackstate = true
						if m.player.CurHP >0{
							go m.AttackSequence()
						}else{
							m.selCom+=2
						}
					case "Item":					
						m.attributes = m.UserInv.ItemName
						InventoryView(m.cfgState, m.UserInv,"item",m.selitem)					
						m.output= fmt.Sprintf("Switching from the %d available ones,at %d",len(m.attributes),m.selitem)

					case "Switch":
						m.player = nil								
						m.attributes = m.UserInv.PokeName										
						InventoryView(m.cfgState,m.UserInv,"switch",m.selswitch)
						m.player = InitPlpokeStats(m.UserInv.PokeDescriptions)
						m.output = fmt.Sprintf("Switching from the %d available ones,at %d",len(m.attributes),m.selswitch)
					case "Catch":
						m.output = "Catching Pokemon, choose the balls!"
					case "ViewStats":
						m.output ="Displaying stats"
					case "Escape":
						m.output = "You escaped the battle!"
						m.battlestate = false
					}
					m.textInput.SetValue("")

				// }else if m.battlestate && m.attackstate{
				// 	var endev string
				// 	m.textInput.SetValue(m.UserInv.MoveDescriptions[m.selmove])
				// 	dev:= playerAttack(m.cfgState,m.textInput.Value(),m.MoAt,m.player,m.enemy)
				// 	m.output = fmt.Sprintf("%s,%d,%d,%s,%d,%s",m.textInput.Value(),m.player.MaxHP,m.player.Attack,m.UserInv.PokeName[m.selswitch],m.MoAt.endamage,dev)
				// 	if dev == "Success"{
				// 		endev = enemyAttack(m.cfgState,m.enemy,m.player,m.MoAt)
				// 	}
				// 	if endev == dev{
				// 		m.output = "Turn complete"
				// 		m.attackstate = false
				// 	}else{
				// 		m.output = "Turn complete with missed attacks"
				// 		m.attackstate = false
				// 	}
					
					
				}else {
					m.output, m.locationList.Items = processCommand(m.textInput.Value(), m.cfgState)
					m.locationList.count = len(m.locationList.Items)
					m.showPoke =false
					m.showLoc = true
					m.battlestate = false
					m.showmsg = false
					m.attackstate = false
					m.textInput.SetValue("")
				}		
			}
		case clearMessage:
			// Clears the initial battle message and switch to the battle HUD
			m.showmsg = false
			m.output = ""
			return m, cmd
		}

	m.textInput, cmd = m.textInput.Update(msg)
	}
	return m,cmd
}

func InventoryView(cfgState *config.ConfigState,uvin *UserInventory,swcase string, index int){
	if len(cfgState.PokemonCaught) == 0 && len(cfgState.ItemsHeld) == 0{
		return 
	}	
	switch swcase{
	case "fill":
		uvin.PokeName=nil
		uvin.ItemName =nil
		for _,poke := range cfgState.PokemonCaught{
			uvin.PokeName = append(uvin.PokeName, poke.Name)
		}
		for _,item := range cfgState.ItemsHeld{
			uvin.ItemName = append(uvin.ItemName, item.Name)
		}
	case "switch":
		uvin.PokeDescriptions =nil
		for _, poke := range cfgState.PokemonCaught{
			if poke.Name == uvin.PokeName[index]{
				ascii_img, _ := imagegen.AsciiGen(poke.Sprites.FrontDefault,52)
				uvin.PokeSprite =  ascii_img
				uvin.PokeDescriptions = append(uvin.PokeDescriptions, "LV1")
				for _,stats := range poke.Stats{
					uvin.PokeDescriptions = append(uvin.PokeDescriptions, strconv.Itoa(stats.BaseStat))
				}
				uvin.PokeDescriptions = append(uvin.PokeDescriptions, uvin.PokeDescriptions[1])
			}
		}	
	case "item":
		for _, item := range cfgState.ItemsHeld{
			uvin.ItemDescriptions= append(uvin.ItemDescriptions, item.Category.Name)
			uvin.ItemDescriptions = append(uvin.ItemDescriptions, item.EffectEntries[0].ShortEffect)
		}
	case "move":
		uvin.MoveName=nil
		uvin.MoveDescriptions = nil
		for _, poke := range cfgState.PokemonCaught{
			if poke.Name == uvin.PokeName[index]{
				for i, move := range poke.Moves{
					if i<3{
						uvin.MoveName = append(uvin.MoveName, move.Move.Name)
						uvin.MoveDescriptions = append(uvin.MoveDescriptions, move.Move.URL)
					}
				}
			}}

	}	
}

func(m *btBaseModel) AttackSequence(){
	m.moveMutex.Lock()
	defer m.moveMutex.Unlock()
	m.textInput.SetValue(m.UserInv.MoveDescriptions[m.selmove])
	plres := playerAttack(m.cfgState,m.textInput.Value(),m.MoAt,m.player,m.enemy)
	m.output += fmt.Sprintf("Player's attack result: %s", plres)
	if m.enemy.CurHP >0{
		enres := enemyAttack(m.cfgState,m.enemy,m.player,m.MoAt)
		m.output += fmt.Sprintf("Player's attack result: %s", enres)
	}
	m.turnComplete <-true
}


func playerAttack(cfgState *config.ConfigState,MoveURL string,movat *MoveAttack, plpokstat *PlayerPokemonStats, enpokstat *EnemyPokemonStats) string{
	
	moveStats, err := cfgState.PokeapiClient.InvokeMove(MoveURL)
	if err!=nil { return  "Move values not imported"}
	hitChance := rng.Intn(12)// 1/12 chance to miss
	hitmult := rng.Float64()
	if moveStats.Accuracy <56{
		return "low Accuracy"
	}
	
	if moveStats.DamageClass.Name == "physical" {
		if hitChance == 0 {
			return "physical attack failed"
		}
		enemy_defense := enpokstat.Defense	
		player_damage := plpokstat.Attack 
		basedamage := int(float64((moveStats.Power *(player_damage/enemy_defense))+20)*hitmult)
		if moveStats.Meta.CritRate > 1 && rng.Intn(16) < moveStats.Meta.CritRate{
			movat.endamage = int(float64(basedamage)*1.5) 
		}
		movat.endamage = basedamage 
		if moveStats.Meta.AilmentChance > 0 && rand.Intn(100) < moveStats.Meta.AilmentChance {
			movat.enStatusEffect = moveStats.Meta.Ailment.Name
		}
	
	}else if moveStats.DamageClass.Name == "special"{
		if hitChance == 12 || hitChance ==7 {
			return "Special Attack failed"
		}
		sed := enpokstat.SpecialDefense
		spd := plpokstat.SpecialAttack
		bd :=  int(float64((moveStats.Power *(spd/sed))+20)*hitmult)
		if rng.Intn(24) < moveStats.Meta.CritRate{
			movat.endamage = int(float64(bd)*1.5) 
		}
		movat.endamage = bd
		if moveStats.Meta.AilmentChance > 0 && rand.Intn(100) < moveStats.Meta.AilmentChance {
			movat.enStatusEffect = moveStats.Meta.Ailment.Name
		}

	}else if  moveStats.DamageClass.Name == "status"{
		sed := enpokstat.SpecialDefense
		spd := enpokstat.SpecialAttack
		bd :=  int(float64((moveStats.Power *(spd/sed))+20)*hitmult)
		if moveStats.Meta.CritRate>4 && rng.Intn(24) < moveStats.Meta.CritRate{
			movat.endamage = int(float64(bd)*1.5) 
		}
		movat.endamage = bd
		if moveStats.Meta.AilmentChance > 0 {
			movat.enStatusEffect = moveStats.Meta.Ailment.Name}
	}else{
		return "Nothing there in moveStats"
	}
	
	curhp:= enpokstat.CurHP
	curhp = curhp - movat.endamage
	if curhp < 0{
		enpokstat.CurHP = 0
	}
	enpokstat.CurHP = curhp

	return "Success"
}

func randFloat(min, max float64) float64 {
    return min + rng.Float64()*(max-min)
}

func splitString(s string) (string, string) {
    var index int
    for i, r := range s {
        if unicode.IsDigit(r) {
            index = i
            break
        }
    }
    return s[:index], s[index:]
}

func enemyAttack(cfgState *config.ConfigState,enemy *EnemyPokemonStats,player *PlayerPokemonStats,move *MoveAttack) string{
	moveurl := "https://pokeapi.co/api/v2/move/29"
	moveStats, err := cfgState.PokeapiClient.InvokeMove(moveurl)
	if err!=nil{return "Url at Fault"}
	_,lvl := splitString(enemy.Level)
	level, converr := strconv.Atoi(lvl)
	if converr != nil{return "Unable to convert"}
	attackMiss := rng.Intn(36)
	a,b := 0.08, 0.29
	enpower := randFloat(a,b)

	if attackMiss != 29{
		pdef := player.Defense	
		edam := enemy.Attack 
		basedamage := int((moveStats.Power *(edam/pdef))/int(math.Round((float64(level*10) * enpower))))
		if moveStats.Meta.CritRate > 1 && rng.Intn(8) < moveStats.Meta.CritRate{
			move.pldamage = int(float64(basedamage)*1.5) 
		}
		move.pldamage = basedamage 
		if moveStats.Meta.AilmentChance > 0 && rand.Intn(100) < moveStats.Meta.AilmentChance {
			move.plStatusEffect = moveStats.Meta.Ailment.Name
		}
	}

	curhp := player.CurHP
	curhp = curhp - move.pldamage
	if curhp <=0{
		player.CurHP =0
	}else{player.CurHP = curhp}

	return "Success"

}


type clearMessage struct{}

type BattleLog struct{
	logmsg 	 string
	msgarr	 []string
}

type UserInventory struct{
	PokeSprite				string
	PokeName				[]string
	ItemName				[]string
	MoveName				[]string
	PokeDescriptions		[]string
	ItemDescriptions		[]string
	MoveDescriptions		[]string
}

type MoveAttack struct{
	enStatusEffect	string
	plStatusEffect	string
	endamage			int
	pldamage			int
}

type PlayerPokemonStats struct{
	Level		   string
	MaxHP          int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
	CurHP		   int
}

type EnemyPokemonStats struct{
	Name		   string
	Level		   string
	MaxHP          int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
	CurHP		   int
	BaseExperience int
	Sprites		   string	
}

func InitInventoryListing() *UserInventory {
	return &UserInventory{
		PokeSprite:  	  "",
		PokeName: 		  []string{},
		ItemName: 		  []string{},
		MoveName: 		  make([]string, 0, 3),	
		PokeDescriptions: []string{},
		ItemDescriptions: []string{},
		MoveDescriptions: make([]string, 0, 3),
	}
}

func InitBattlelogIniate() *BattleLog{
	return &BattleLog{
		logmsg:   "",
		msgarr:   make([]string,40,200),
	}
}

func InitMovatt() *MoveAttack{
	return &MoveAttack{
		enStatusEffect: "",
		plStatusEffect: "",
		endamage: 		0,
		pldamage: 		0,
	}
}

func InitPlpokeStats(values []string) *PlayerPokemonStats{
	return &PlayerPokemonStats{
		Level: values[0],
		MaxHP: 			parseInt(values[1]),
		Attack: 		parseInt(values[2]),
		Defense:        parseInt(values[3]),
		SpecialAttack:  parseInt(values[4]),
		SpecialDefense: parseInt(values[5]),
		Speed:          parseInt(values[6]),
		CurHP:			parseInt(values[1]),		
	}
}

func InitEnpokeStats(values []string) *EnemyPokemonStats {
	return &EnemyPokemonStats{
		Name:			values[0],
		Level:          values[1],
		MaxHP:          parseInt(values[2]),
		Attack:         parseInt(values[3]),
		Defense:        parseInt(values[4]),
		SpecialAttack:  parseInt(values[5]),
		SpecialDefense: parseInt(values[6]),
		Speed:          parseInt(values[7]),
		CurHP:          parseInt(values[8]),
		BaseExperience: parseInt(values[9]),
		Sprites: 		values[10],
	}
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func HealthBar(curHP,maxHP int)string{

	hpercent := float64(curHP)/float64(maxHP)
	barwidth := 10

	filledWidth :=int(hpercent*float64(barwidth))
	if filledWidth < 0 {
		filledWidth = 0
	}else if filledWidth > barwidth { filledWidth = barwidth}
	
	emptyWidth:= barwidth - filledWidth
	if emptyWidth < 0 {
		emptyWidth = 0
	}
	
	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#004400"))
    emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Background(lipgloss.Color("#222222"))

	filled := filledStyle.Render(strings.Repeat("█", filledWidth))
	empty :=  emptyStyle.Render(strings.Repeat(" ",emptyWidth))
	hpercent *= 100
	return fmt.Sprintf("%s%s %d%%", filled, empty, int(hpercent))

}

func (m *btBaseModel) View() string {
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
		visibleCommands := m.commands[m.trackcomind:min(m.trackcomind+cmdview,len(m.commands))]
		viscom := make([]string, cmdview)
		for i, cmd := range visibleCommands {
			style := lipgloss.NewStyle().Bold(true)
			if m.trackcomind+i == m.selCom {
				style = style.Foreground(lipgloss.Color("220")).Underline(true)
			} else {
				style = style.Foreground(lipgloss.Color("240"))
			}
			viscom[i] = style.Render("[" + cmd + "]")
		}
		attributes := make([]string, len(m.attributes))
		for i,attr := range m.attributes{
			style := lipgloss.NewStyle().Bold(true)
			if (m.commands[m.selCom] == "Attack" && i == m.selmove) ||
			(m.commands[m.selCom] == "Item" && i == m.selitem) ||
			(m.commands[m.selCom] == "Switch" && i == m.selswitch){
				style = style.Foreground(lipgloss.Color("220")).Underline(true)
			}else{
				style = style.Foreground(lipgloss.Color("240"))
			}
			attributes[i] = style.Render("[" + attr+ "]")
		}
		var plstats string
		var tlstats string

		commandBoxContent := lipgloss.JoinVertical(lipgloss.Top, viscom...)
		commandAttrContent := lipgloss.JoinHorizontal(lipgloss.Top, attributes...)

		emptybox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).BorderForeground(lipgloss.Color("#008080")).Padding(0, 0).Width(10).Height(5).Align(lipgloss.Left).Render()
		commandAttributes:= lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#008080")).Padding(0, 0).Width(37).Height(5).Align(lipgloss.Left).Render(commandAttrContent)
		commandListBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("220")).Padding(0, 0).Width(20).Height(5).Align(lipgloss.Left).Render(commandBoxContent)
		battlelog := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("220")).Padding(0,0).Width(81).Height(5).Align(lipgloss.Center).Render("Battle Log")

		
		commandBox := lipgloss.JoinHorizontal(lipgloss.Left,emptybox,commandListBox,
			commandAttributes,battlelog)

		
		enemyViewBox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0,0,0,0).Width(64).Height(15).Align(lipgloss.Center)
		enemyStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#00CFFF")).Padding(0,0,0,0).Width(40).Height(3).Align(lipgloss.Center)
		
		playerViewBox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0, 0, 0).Width(64).Height(15).Align(lipgloss.Center)
		playerStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#F25D94")).Padding(0,0,0,0).Width(40).Height(3).Align(lipgloss.Center)
		

		//enemy stats box
		
		maxhp := m.enemy.MaxHP
		curhp := m.enemy.CurHP
		enemyhb := HealthBar(curhp,maxhp)
		elvlhb := lipgloss.JoinHorizontal(lipgloss.Left,m.enemy.Level," ",enemyhb)
		tlstats = lipgloss.JoinVertical(lipgloss.Top,cases.Title(language.Und, cases.NoLower).String((m.enemy.Name)),elvlhb)
				
		// player stats box
		if 	m.player.MaxHP !=0 {
			maxhp := m.player.MaxHP
			curhp := m.player.CurHP
			playerhb := HealthBar(curhp, maxhp)
			plvlhb := lipgloss.JoinHorizontal(lipgloss.Left,m.player.Level," ",playerhb)
			plstats = lipgloss.JoinVertical(lipgloss.Top, cases.Title(language.Und, cases.NoLower).String((m.UserInv.PokeName[m.selswitch])), plvlhb)
		} else {
			plstats = lipgloss.JoinVertical(lipgloss.Top, "Player 1")
		}

		topLeftPokemon := lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(m.enemy.Sprites)
		topLeftStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(tlstats)
		
		bottomRightPokemon := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CFFF")).Render(m.UserInv.PokeSprite)
		bottomRightStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CFFF")).Render(plstats)

		
		// Display Pokémon and commands
		enemyPokemonView := lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.PlaceHorizontal(0,lipgloss.Left,topLeftPokemon))
		
		playerPokemonView := lipgloss.JoinVertical(lipgloss.Bottom,
			lipgloss.PlaceHorizontal(64,lipgloss.Right,bottomRightPokemon))
		
		enemyBox := lipgloss.JoinVertical(lipgloss.Center,enemyViewBox.Render(enemyPokemonView),
					playerStatusBox.Render(topLeftStatus))
		playerBox := lipgloss.JoinVertical(lipgloss.Center, enemyStatusBox.Render(bottomRightStatus),
					 playerViewBox.Render(playerPokemonView))			
		
		
		battleBox :=lipgloss.JoinHorizontal(lipgloss.Left,enemyBox,
					playerBox)


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
