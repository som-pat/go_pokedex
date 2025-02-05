package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"unicode"
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
	"github.com/som-pat/poke_dex/app"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
const cmdview = 4

type btBaseModel struct {
	textInput    	textinput.Model
	output       	string
	cfgState     	*config.ConfigState
	locationList 	*Paginatedlisting
	PokemonList  	*Paginatedlisting
	UserInv			*UserInventory
	MoAt			*MoveAttack
	enemy			*EnemyPokemonStats
	player			*PlayerPokemonStats
	AttackEvents	[]AttackEvent
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
	logchannel   	chan string
	battlelog       []string
	showinspect		bool
	// checkplshake    bool
	// plshakech		chan bool
	// checkenshake	bool
	eventchan	    chan bool
	comsel 			bool
	width 			int
	height 			int
	curevent		int
	Navigator		*app.AppNavigator
}

type Paginatedlisting struct {
	Items        []string
	count        int
	selectedIndex int
}

type LogMessage struct{
	log 	string
}

func takeInput(cfgState *config.ConfigState, navigator *app.AppNavigator) *btBaseModel {
	ti := textinput.New()
	ti.Placeholder = "[Press 'esc' to quit, 'up'/'down'/'right'/'left' to navigate, 'enter' to confirm]"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 200
	return &btBaseModel{
		textInput:    ti,
		output:       "Welcome to PokeCLI!\nType 'explore' to search for Pokémon, 'catch' to catch one, or 'quit' to exit.",
		cfgState:     cfgState,
		locationList: InitPaginatedListing(20),
		PokemonList:  InitPaginatedListing(40),
		UserInv: 	  &UserInventory{},
		width:		  0,
		height:       0,
		enemy: 		  &EnemyPokemonStats{},
		player:		  &PlayerPokemonStats{},
		AttackEvents: make([]AttackEvent, 0),	
		showLoc: 	  false,
		showPoke: 	  false,
		battlestate:  false,
		battleTarget: "",
		showmsg:      false,
		selCom: 	  0,
		commands: 	  []string{"Attack","Item","Switch","Catch","Stats","Escape"},
		attributes:   []string{},
		selswitch: 	  0,
		selmove: 	  0,
		selitem: 	  0,
		trackcomind:  0,
		attackstate:  false,
		MoAt: 		  &MoveAttack{},
		moveMutex:	  sync.Mutex{},	
		turnComplete: make(chan bool,1),
		logchannel:	  make(chan string,4),
		battlelog:    []string{},
		showinspect:  false,
		// plshakech:    make(chan bool,1),
		// checkplshake: false,
		// checkenshake: false,
		comsel:       false,
		curevent: 	  0,
		eventchan:	  make(chan bool,1),
		Navigator: 	  navigator,	 
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

func (m *btBaseModel) startLogging(){
	go m.processLogs()
}

func (m *btBaseModel) processLogs(){
	for logMsg := range m.logchannel {
		if m.battlestate { // Log messages only if in battle state
			m.battlelog = append(m.battlelog, logMsg)
			if len(m.battlelog) > 4 {
				m.battlelog = m.battlelog[len(m.battlelog)-4:] // Keep only the last 5 messages
			}
		}else{
			return
		}
	}
}


func (m *btBaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	select {
	case <-m.turnComplete:
		m.attackstate = false
		m.logchannel <-"Turn complete. Choose your next move."

	case logMsg := <-m.logchannel:
		if m.battlestate{
		return m, func() tea.Msg {
			return LogMessage{log: logMsg}
			}
		}

	case event := <- m.eventchan:
		if event {
			cmd = m.eventsequence()
			return m,cmd
		} 
		


	default:
		switch msg := msg.(type) {		

		case LogMessage:
			return m,cmd
		case tea.WindowSizeMsg:
			m.width  = msg.Width
			m.height = msg.Height
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				return m.Navigator.GoToMenu(), nil //battlemode m.Navigator.GoToMenu()
			case "left":
				switch m.commands[m.selCom]  {
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
					m.PokemonList.count = len(m.PokemonList.Items)
					m.showinspect = false
					m.showPoke =true
					m.showLoc = false
					m.battlestate = false
					m.showmsg = false
					m.attackstate = false
					
					m.textInput.SetValue("") 

				}else if strings.HasPrefix(m.textInput.Value(),"battle") {
					m.attributes = nil
					m.selmove = 0 
					m.selCom = 0
					m.selitem = 0
					m.selswitch = 0
					if m.UserInv == nil{
						m.UserInv = InitInventoryListing()					
					}
					InventoryView(m.cfgState,m.UserInv,"fill",0)
					m.battlestate = true
					m.showLoc = false
					m.showPoke = true
					m.showmsg = true
					m.attackstate = false
					m.showinspect = false
					m.output,m.PokemonList.Items = processCommand(m.textInput.Value(),m.cfgState)
					if m.enemy != nil{m.enemy = nil}
					m.enemy = InitEnpokeStats(m.PokemonList.Items)
					m.battleTarget = m.enemy.Name
					m.PokemonList.Items = nil			
					m.textInput.SetValue("")
					return m, tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
						return clearMessage{}
					})

				}else if m.battlestate && !m.attackstate && !m.comsel{
					m.comsel = true
					switch m.commands[m.selCom] {
					case "Attack":
						m.UserInv.MoveName =nil
						InventoryView(m.cfgState,m.UserInv,"move",m.selswitch)
						m.attributes = m.UserInv.MoveName
					case "Item":
						m.attributes = m.UserInv.ItemName
					case "Switch":
						m.attributes = m.UserInv.PokeName
					case "Catch":
						m.attributes = []string{}
					case "ViewStats":
						m.attributes = []string{}
					case "Escape":
						m.handleEscape()
					}
				}else if m.battlestate && !m.attackstate && m.comsel{
					m.comsel = false
					if !m.comsel {}
					m.battlelog =nil
					m.startLogging()				
					switch m.commands[m.selCom] {
					case "Attack":				
						m.UserInv.MoveName =nil
						InventoryView(m.cfgState,m.UserInv,"move",m.selswitch)
						m.attributes = m.UserInv.MoveName
						m.attackstate = true
						m.AttackEvents = nil
						if m.player.CurHP >1 && m.enemy.CurHP>1{							
							go m.AttackSequence()
						}else{
							m.textInput.SetValue("")
							m.handleSwitch()
						}
					case "Item":					
						m.attributes = m.UserInv.ItemName
						InventoryView(m.cfgState, m.UserInv,"item",m.selitem)					
						m.logchannel<-fmt.Sprintf("Using item %d out of %d",m.selitem,len(m.attributes))

					case "Switch":
						m.handleSwitch()

					case "Catch":
						m.logchannel <- "Catching Pokemon, choose the balls!"
					case "Stats":
						m.logchannel <- "Displaying stats"
						m.logchannel <- fmt.Sprintf("Lvl: %s, MaxHP: %d, Att: %d, Def: %d",m.player.Level,m.player.MaxHP,m.player.Attack,m.player.Defense)
						m.logchannel <- fmt.Sprintf("SpAtk: %d, SpDef: %d, Spd: %d",m.player.SpecialAttack,m.player.SpecialDefense,m.player.Speed)
					case "Escape":
						m.handleEscape()
					}
					m.textInput.SetValue("")
					return m, cmd
					
					
				}else {
					m.output, m.locationList.Items = processCommand(m.textInput.Value(), m.cfgState)
					m.locationList.count = len(m.locationList.Items)
					m.showinspect =false
					part := strings.Split(m.textInput.Value(), " ")
					if part[0] =="inspect"{m.showinspect = true}
					m.showLoc = true
					m.showPoke =false				
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

func (m *btBaseModel) handleSwitch() {
	m.player = nil
	if len(m.UserInv.PokeName)>0{
		m.attributes = m.UserInv.PokeName
		InventoryView(m.cfgState, m.UserInv, "switch", m.selswitch)
		m.player = InitPlpokeStats(m.UserInv.PokeDescriptions)
		m.logchannel <- fmt.Sprintf("Switching to %d from available %d",m.selswitch, len(m.attributes))
	}else{
		m.logchannel <- "No pokemon to switch with"
		m.textInput.SetValue("")
	}
}



func (m *btBaseModel) handleEscape() {
	m.logchannel <- "Escaping in dire straits, I see"
	time.Sleep(1*time.Second)
	m.battlestate = false
	m.battlelog = nil
	m.attackstate = false
	m.output = "You escaped the battle!"
	m.textInput.SetValue("")
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
				ascii_img, _ := imagegen.AsciiGen(poke.Sprites.BackDefault,52)
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

type Battle struct{
	AttackEvents	[]AttackEvent
	Defeated		[]string
	UserInventory	*UserInventory
}


type AttackEvent struct {
	EventType 		string
	attackSprite	string
	playershake		bool
	enemyshake		bool
	showatsprite	bool

}

func InitAtevent()	AttackEvent{
	return AttackEvent{
		EventType: 		"",
		attackSprite: 	imagegen.AttackGen("imagegen/sprites/2.png"),
		playershake: 	false,
		enemyshake: 	false,
		showatsprite: 	false,
	}
}

func(m *btBaseModel) AttackSequence(){
	enchance  := rng.Intn(5)
	m.moveMutex.Lock()
	defer m.moveMutex.Unlock()
	if enchance == 2{
		m.Enemyfirst()
		time.Tick(2*time.Second)		
		if m.player.CurHP >1{
			m.PlayerFirst()
		}else {
			m.player.CurHP = 0 
			m.handleSwitch()}
	}else{
		m.PlayerFirst()
		time.Sleep(2*time.Second)	
		if m.enemy.CurHP >1{
			m.Enemyfirst()
		}else{
			m.enemy.CurHP = 0
			m.logchannel<-fmt.Sprintf("Successfully Defeated, Received %d",m.enemy.BaseExperience)
		}
	}
	
	// m.logchannel <-"Player attacks.."
	// m.textInput.SetValue(m.UserInv.MoveDescriptions[m.selmove])
	// plres,dam := playerAttack(m.cfgState,m.textInput.Value(),m.MoAt,m.player,m.enemy)
	// m.logchannel <- fmt.Sprintf("Player's attack result: %d ", dam)
	// if dam >0{m.enshakech <- true}
	// if plres !="Success" || m.enemy.CurHP <=1{
	// 	m.logchannel<-fmt.Sprintf("Successfully Defeated, Received %d",m.enemy.BaseExperience)
	// }else{
	// 	time.Sleep(800 *time.Microsecond)
	// 	m.logchannel <- "Enemy attacks..."
	// 	_,dam := enemyAttack(m.cfgState,m.enemy,m.player,m.MoAt)
	// 	if dam >0 {m.plshakech <- true}		
	// 	m.logchannel <- fmt.Sprintf("Enemy's attack result: %d", dam)
	// }
	m.turnComplete <-true

}

func (m *btBaseModel) Enemyfirst() {
	m.logchannel <- "Enemy attacks..."
	_,dam := enemyAttack(m.cfgState,m.enemy,m.player,m.MoAt)
	if dam >0 {	
		m.curevent = 0	
		m.AttackEvents = []AttackEvent{
		{EventType: "SlashAnimation",attackSprite:imagegen.AttackGen("imagegen/sprites/attack/2.png"),
		playershake:true,enemyshake:false,showatsprite:true},

		{EventType: "CameraShake",attackSprite:"",
		playershake:true,enemyshake: false,showatsprite: false},
		}
	for i:=0;i<2;i++{
		m.curevent = i			
		m.eventchan <-true
	}
	

	}		
	m.logchannel <- fmt.Sprintf("Enemy's attack result: %d", dam)
}

func (m *btBaseModel) PlayerFirst(){
	m.logchannel <-"Player attacks.."
	m.textInput.SetValue(m.UserInv.MoveDescriptions[m.selmove])
	_,dam := playerAttack(m.cfgState,m.textInput.Value(),m.MoAt,m.player,m.enemy)

	if dam >0 {
		m.curevent = 0
		m.AttackEvents = []AttackEvent{
			{EventType: "SlashAnimation",attackSprite:imagegen.AttackGen("imagegen/sprites/attack/2.png"),
			playershake:false,enemyshake:true,showatsprite:true},
			
			{EventType: "CameraShake",attackSprite:"",
			playershake:false,enemyshake: true,showatsprite: false},	
			}

		for i:=0;i<2;i++{			
			m.eventchan <-true
		}
		
	}
	m.logchannel <- fmt.Sprintf("Player's attack result: %d ", dam)
}

func (m *btBaseModel) eventsequence() tea.Cmd{
	if m.curevent >= len(m.AttackEvents){
		m.AttackEvents = nil
		return nil
	}

	CurrentEvent := m.AttackEvents[m.curevent]
	switch CurrentEvent.EventType{

	case "SlashAnimation":
		m.output+= "0.3  "
		return tea.Tick(4*time.Second, func(t time.Time) tea.Msg { 
			return t
		})
	case "CameraShake":
		m.output += "1.3  "
		if CurrentEvent.playershake{			
			return tea.Tick(3*time.Second, func(t time.Time) tea.Msg { 
				return t
			})}
		if CurrentEvent.enemyshake{			
			return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
				CurrentEvent.playershake = false 
				return t
			})}
		
		}

	return nil		
}

func playerAttack(cfgState *config.ConfigState,MoveURL string,movat *MoveAttack, plpokstat *PlayerPokemonStats, enpokstat *EnemyPokemonStats) (string,int){
	
	moveStats, err := cfgState.PokeapiClient.InvokeMove(MoveURL)
	if err!=nil { return  "Move values not imported",0}
	hitChance := rng.Intn(12)// 1/12 chance to miss
	hitmult := rng.Float64()
	if moveStats.Accuracy <56{
		return "low Accuracy",0
	}
	
	if moveStats.DamageClass.Name == "physical" {
		if hitChance == 0 {
			return "physical attack failed",0
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
			return "Special Attack failed",0
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
		return "Nothing there in moveStats",0
	}
	
	curhp:= enpokstat.CurHP
	curhp = curhp - movat.endamage
	if curhp < 0{
		enpokstat.CurHP = 0
	}
	enpokstat.CurHP = curhp

	return "Success",movat.endamage
}

func randFloat(min, max float64) float64 {
    return min + rng.Float64()*(max-min)
}

func DigsplitString(s string) (string, string) {
    var index int
    for i, r := range s {
        if unicode.IsDigit(r) {
            index = i
            break
        }
    }
    return s[:index], s[index:]
}

func enemyAttack(cfgState *config.ConfigState,enemy *EnemyPokemonStats,player *PlayerPokemonStats,move *MoveAttack) (string, int){
	moveurl := "https://pokeapi.co/api/v2/move/29"
	moveStats, err := cfgState.PokeapiClient.InvokeMove(moveurl)
	if err!=nil{return "Url at Fault",0}
	_,lvl := DigsplitString(enemy.Level)
	level, converr := strconv.Atoi(lvl)
	if converr != nil{return "Unable to convert",0}
	attackMiss := rng.Intn(512)
	a,b := 0.08, 0.29
	enpower := randFloat(a,b)

	if attackMiss != 126{
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
	}else{
		return "Attack Missed",0
	}

	curhp := player.CurHP
	curhp = curhp - move.pldamage
	if curhp <=0{
		player.CurHP =0
	}else{player.CurHP = curhp}

	return "Success",move.pldamage

}

type clearMessage struct{}

// type BattleLog struct{
// 	logmsg 	 string
// 	msgarr	 []string
// }

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
	endamage		int
	pldamage		int
	plspritelap		string
	enspritelap     string

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
	spriteurl	   string	
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

// func InitBattlelogIniate() *BattleLog{
// 	return &BattleLog{
// 		logmsg:   "",
// 		msgarr:   make([]string,40,200),
// 	}
// }

func InitMovatt() *MoveAttack{
	return &MoveAttack{
		enStatusEffect: "",
		plStatusEffect: "",
		endamage: 		0,
		pldamage: 		0,
		enspritelap: 	"",
		plspritelap:    "",	
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
		spriteurl: 		values[11],
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
	if curHP <=0 {curHP =0}

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



func DamageShake(frac float64)(int, int){
	if frac<=0 || frac >=1{
		return 0,0
	}
	mid:= 0.25

	offsetx := 4.5 * math.Sin(20.8 * math.Pi + 2.3)+
			   2.1 * math.Sin(5.6 * math.Pi + 5)+
			   1.9 * math.Sin(11.2 *math.Pi +1.2)
	offsety	:= 4 * math.Sin(16.8 * math.Pi - 2.3)+
			   2.8 * math.Sin(8.6 * math.Pi + 4)+
			   1.2 * math.Sin(1.2 *math.Pi +0)
	
	var lim float64
	if frac <= mid{
		lim = (1/mid) * frac
	}else{
		lim = 1- (frac-mid)/(1-mid)
	}
	return int(offsetx * lim),int(offsety * lim)

}


// func overlay(base, toverlay string, offsetx, offsety int) string {
//     baselines := strings.Split(base, "\n")
//     overlines := strings.Split(toverlay, "\n")

//     for y, line := range overlines {
//         if y+offsety >= len(baselines) {
//             continue
//         }

//         baseLine := []rune(baselines[y+offsety])

//         for x, char := range line {
//             if x+offsetx < len(baseLine) && char != ' ' {
//                 baseLine[x+offsetx] = char
//             }
//         }
//         baselines[y+offsety] = string(baseLine)
//     }

//     return strings.Join(baselines, "\n")
// }



func (m *btBaseModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("35"))
	screenHeight := 44	
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
		var enstats string


		commandBoxContent := lipgloss.JoinVertical(lipgloss.Top, viscom...)
		commandAttrContent := lipgloss.JoinHorizontal(lipgloss.Top, attributes...)
		logContent := lipgloss.JoinVertical(lipgloss.Left, m.battlelog...)

		// emptybox := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0).Width(10).Height(4).Align(lipgloss.Left).Render()
		commandAttributes:= lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("220")).
									Padding(0, 0).Width(38).Height(4).Align(lipgloss.Left).Render(commandAttrContent)
		commandListBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("220")).
									Padding(0, 0).Width(18).Height(4).Align(lipgloss.Center).Render(commandBoxContent)
		battlelog := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#008080")).
									Padding(0,0).Width(66).Height(4).Align(lipgloss.Left).Render(logContent)

		
		commandBox := lipgloss.JoinHorizontal(lipgloss.Left,commandListBox,
			commandAttributes,battlelog)
	
		
		enpokebox := lipgloss.NewStyle().Border(lipgloss.NormalBorder(),false,false,false,false).
		Padding(0,0,0,0).Width(64).Height(10).Align(lipgloss.Center)
		plpokebox := lipgloss.NewStyle().Border(lipgloss.NormalBorder(),false,false,false,false).
		Padding(0,0,0,0).Width(64).Height(10).Align(lipgloss.Center) 
		
		enstatbox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#F25D94")).
								Padding(0,0,0,0).Width(40).Height(2).Align(lipgloss.Center)
		plstatbox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#00CFFF")).
								Padding(0,0,0,0).Width(40).Height(2).Align(lipgloss.Center)
		

		roundedbg := imagegen.BgMaker("imagegen/sprites/bg/rectangle1.png") 	
		enbg := lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32")).Border(lipgloss.NormalBorder(),false,false,false,false).
				Padding(0,0,4,0).Width(60).Height(1).Align(lipgloss.Center)
		plbg := lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32")).Border(lipgloss.NormalBorder(),false,false,false,false).
		Padding(0,0,0,0).Width(60).Height(1).Align(lipgloss.Center)
		

	// vieb := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0).Width(10).Height(10).Align(lipgloss.Left).Render()
	// 	// vipl := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0).Width(10).Height(10).Align(lipgloss.Left).Render()
	// 	enemyViewBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0).Width(fixedBoxWidth).Height(fixedBoxHeight).Align(lipgloss.Center)
	// 	enemyStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#F25D94")).
	// 						Padding(0,0,0,0).Width(40).Height(2).Align(lipgloss.Center)
		
	// 	playerViewBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0).Width(fixedBoxWidth).Height(fixedBoxHeight).Align(lipgloss.Center)
	// 	playerStatusBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#00CFFF")).
	// 						Padding(0,0,0,0).Width(40).Height(2).Align(lipgloss.Center)
		

		//enemy stats box
		
		maxhp := m.enemy.MaxHP
		curhp := m.enemy.CurHP
		enemyhb := HealthBar(curhp,maxhp)
		elvlhb := lipgloss.JoinHorizontal(lipgloss.Left,m.enemy.Level," ",enemyhb)
		enstats = lipgloss.JoinVertical(lipgloss.Top,cases.Title(language.Und, cases.NoLower).String((m.enemy.Name)),elvlhb)
				
			// player stats box
		if 	m.player.MaxHP !=0 {
			maxhp := m.player.MaxHP
			curhp := m.player.CurHP
			playerhb := HealthBar(curhp, maxhp)
			plvlhb := lipgloss.JoinHorizontal(lipgloss.Left,m.player.Level," ",playerhb)
			plstats = lipgloss.JoinVertical(lipgloss.Top, cases.Title(language.Und, cases.NoLower).
						String((m.UserInv.PokeName[m.selswitch])), plvlhb)
		} else {
			plstats = lipgloss.JoinVertical(lipgloss.Top, "Player 1")
		}

			
		newh := 52
		for{
			w,_ := lipgloss.Size(m.enemy.Sprites)
			if  w<53{
				newh+=1 
				m.enemy.Sprites,_= imagegen.AsciiGen(m.enemy.spriteurl,newh)
			}else{
				break
			}
					
		}
		var topRightPokemon string
		var bottomLeftPokemon string
		shakexPlayer, shakeyPlayer := 0, 0
    	shakexEnemy, shakeyEnemy := 0, 0
		
		_,h:= lipgloss.Size(m.enemy.Sprites)
		// m.output+= fmt.Sprintf("enem,wid:%d,hei%d \t",w,h)
		if h>=20{
			enpokebox= enpokebox.UnsetPaddingLeft()
			enbg = enbg.UnsetPaddingBottom()
		} 
		topRightPokemon = lipgloss.NewStyle().Render(m.enemy.Sprites)
		bottomLeftPokemon = lipgloss.NewStyle().Render(m.UserInv.PokeSprite) // add to m.player
	// 	// w,h = lipgloss.Size(m.UserInv.PokeSprite)
	// 	// m.output+= fmt.Sprintf("player,wid:%d,hei%d \t",w,h)

		if m.attackstate && m.AttackEvents!=nil{
			if m.AttackEvents[m.curevent].showatsprite{			
				if m.AttackEvents[m.curevent].enemyshake{
					// w,h := lipgloss.Size(m.enemy.Sprites)
					// overlayascii := overlay(m.enemy.Sprites, m.AttackEvents[0].attackSprite, int(w/2), int(h/2))
					topRightPokemon = lipgloss.NewStyle().Render(m.AttackEvents[m.curevent].attackSprite)
					w,h := lipgloss.Size(m.AttackEvents[m.curevent].attackSprite)
					m.output += "0.4ptoe " + fmt.Sprintf("enem,wid:%d,hei%d \t",w,h)
					
				}else if m.AttackEvents[m.curevent].playershake{
					//w,h := lipgloss.Size(m.UserInv.PokeSprite)
					//overlayascii := overlay(m.UserInv.PokeSprite, m.AttackEvents[0].attackSprite, int(w/2), int(h/2))
					bottomLeftPokemon = lipgloss.NewStyle().Render(m.AttackEvents[m.curevent].attackSprite)
					m.output += " 0.5 etop" 
				}
			}else if !m.AttackEvents[m.curevent].showatsprite{
				topRightPokemon = lipgloss.NewStyle().Render(m.enemy.Sprites)
				bottomLeftPokemon = lipgloss.NewStyle().Render(m.UserInv.PokeSprite)
			}
		}



		topRightStatus := lipgloss.NewStyle().Render(enstats)
		bottomLeftStatus := lipgloss.NewStyle().Render(plstats)


		if m.attackstate && m.AttackEvents!=nil{
			if m.AttackEvents[m.curevent].enemyshake && !m.AttackEvents[m.curevent].showatsprite{
				frac:= randFloat(0.2,0.9) 
				shakexEnemy, shakeyEnemy = DamageShake(frac*0.05)
				m.output += "1.3en "
				}
			if m.AttackEvents[m.curevent].playershake && !m.AttackEvents[m.curevent].showatsprite{
				frac := randFloat(0.2, 0.9)
				shakexPlayer, shakeyPlayer = DamageShake(frac*0.05)
				m.output += "1.3pl "
				}
		}		
		entrp := lipgloss.JoinVertical(0.3,
				lipgloss.PlaceHorizontal(60+shakexEnemy,lipgloss.Left,
				lipgloss.PlaceVertical(10+shakeyEnemy,lipgloss.Top,topRightPokemon)))
		plblp := lipgloss.JoinVertical(0.3,
				lipgloss.PlaceHorizontal(60+shakexPlayer,lipgloss.Right,
				lipgloss.PlaceVertical(10+shakeyPlayer, lipgloss.Top,bottomLeftPokemon)))
		// enemyPokemonView := lipgloss.JoinVertical(lipgloss.Top,
		// 	lipgloss.PlaceHorizontal(60,lipgloss.Right,topRightPokemon))
		// enemyPokemonView := lipgloss.JoinVertical(lipgloss.Top,topRightPokemon)
				
		
	// 	playerPokemonView := lipgloss.JoinVertical(lipgloss.Top,
	// 			lipgloss.PlaceHorizontal(60,lipgloss.Right,bottomLeftPokemon))
		
	// 	//posistional box-player
		// posplbox := lipgloss.NewStyle().Border(lipgloss.NormalBorder(),false,false,false,false).
		// 					Padding(0,0,0,0).Width(2).Height(1).Align(lipgloss.Center)
		
		//enpokeattackbox := lipgloss.JoinHorizontal(lipgloss.Left,eattackbox.Render(),enpokebox.Render(topRightPokemon))
		enpokbgbox := lipgloss.JoinVertical(lipgloss.Top,enpokebox.Render(entrp),enbg.Render(roundedbg))
		enemyBox := lipgloss.JoinVertical(lipgloss.Right,enpokbgbox,
					enstatbox.Render(topRightStatus))
		
		//plpokeattackbox := lipgloss.JoinHorizontal(lipgloss.Left,plpokebox.Render(bottomLeftPokemon),pattackbox.Render())
		plpokebgbox := lipgloss.JoinVertical(lipgloss.Top,plblp,plbg.Render(roundedbg))
		playerBox := lipgloss.JoinVertical(lipgloss.Left,plstatbox.Render(bottomLeftStatus), 
					plpokebox.Render(plpokebgbox))			
			
			
		battleBox :=lipgloss.JoinHorizontal(lipgloss.Left,playerBox,enemyBox)


		outer := lipgloss.JoinVertical(lipgloss.Top,battleBox,commandBox)
		outerbox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFFFFF")).Padding(0, 0).Align(lipgloss.Center).Render(outer)
		// lb := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 0).Width(12).MaxHeight(18).Align(lipgloss.Left).Render()
		// rb := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 0).Width(12).MaxHeight(18).Align(lipgloss.Left).Render()
		
							
			centerouterbox := lipgloss.JoinHorizontal(lipgloss.Left,outerbox)
			return lipgloss.JoinVertical(lipgloss.Left,
							centerouterbox,
							m.output,
							
					)
	}else if m.showinspect{
		eb := lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).BorderForeground(lipgloss.Color("220")).
		Padding(0).Width(20).Height(21).Align(lipgloss.Top).Render()
		InspectImage := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("220")).
									Padding(0).Width(90).Height(24).Align(lipgloss.Left).Render(m.locationList.Items[0])
		
		Inspectelem := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("220")).
		Padding(2).Width(28).Height(24).Align(lipgloss.Top).Render(m.output)

		inspectbox := lipgloss.JoinHorizontal(lipgloss.Left, eb ,InspectImage, Inspectelem)
		
		return lipgloss.JoinVertical(lipgloss.Top,
									 inspectbox,
									 "\n",m.textInput.View(),)	
	} 
	

	return fmt.Sprintf(
		"%s\n%s%s %s",
		headerStyle.Render(m.output),
		padding,
		"Pokedex", m.textInput.View(),
	)
}

func processCommand(input string, cfgState *config.ConfigState) (string, []string) {
	input = strings.TrimSpace(input)
	fmt.Printf("Input given %s ", input)
	res, lis := replinternal.ReplInput(cfgState, input)
	return res, lis
}
