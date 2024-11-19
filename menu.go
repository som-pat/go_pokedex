package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/som-pat/poke_dex/imagegen"
	"github.com/som-pat/poke_dex/internal/config"
)

const banner = `
  _____           _        _     __  __                 _            
 |  __ \         | |      | |   |  \/  |               | |           
 | |__) ___   ___| | _____| |_  | \  / | ___  _ __  ___| |_ ___ _ __ 
 |  ___/ _ \ / __| |/ / _ | __| | |\/| |/ _ \| '_ \/ __| __/ _ | '__|
 | |  | (_) | (__|   |  __| |_  | |  | | (_) | | | \__ | ||  __| |   
 |_|   \___/ \___|_|\_\___|\__| |_|  |_|\___/|_| |_|___/\__\___|_|   
`

type Menu struct {
	choices       []string
	cursor        int
	cfgstate      *config.ConfigState
	width         int
	height        int  
	gifFrames     []string
	gifDelay      time.Duration
	currentFrame  int
}
type gifTickMsg time.Time

func MenuModel(cfgstate *config.ConfigState) Menu {
	gifUrl := "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/6.gif"
	frames, delay, err := imagegen.GifGen(gifUrl, 54) 
	if err != nil {
		frames = []string{"Error loading GIF"}
	}
	
	return Menu{
		choices:      []string{"Story Mode", "Battle Mode", "Exit"},
		cursor:       0,
		cfgstate:     cfgstate,
		gifFrames:    frames,
		gifDelay:     delay,
		currentFrame: 0,
	}
}

func (m Menu) Init() tea.Cmd {
  return tea.Tick(m.gifDelay, func(t time.Time) tea.Msg {
		return gifTickMsg(t)
	})
}


func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
  case gifTickMsg:
    m.currentFrame = (m.currentFrame+1)%len(m.gifFrames)
    return m, tea.Tick(m.gifDelay, func(t time.Time) tea.Msg {
      return gifTickMsg(t)
    })

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.choices[m.cursor] {
			case "Battle Mode":
				m.cfgstate.AppState = "BattleMode"
				return takeInput(m.cfgstate), nil
			case "Story Mode":
				m.cfgstate.AppState = "StoryMode"
				return m, nil
			case "Exit":
				return m, tea.Quit
			}
		}
  }
	return m, nil
}

func (m Menu) View() string {
	var options string
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		options += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	gifDisplay := m.gifFrames[m.currentFrame]
	bannerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
	menuStyle := lipgloss.NewStyle().Border(lipgloss.OuterHalfBlockBorder(),false,false,false,false).
		Bold(true).Padding(2).Align(lipgloss.Left)
	menuPlacement := lipgloss.Place(20, 28, lipgloss.Center, lipgloss.Center, menuStyle.Render(options))
	// w,h := lipgloss.Size(gifDisplay)
	// output := fmt.Sprintf("w:%d,h:%d",w,h)
	gifBox := lipgloss.NewStyle().Border(lipgloss.NormalBorder(),false,false,false,false).Width(62).Height(m.height-5).Align(lipgloss.Center).Render(gifDisplay)
	gifPlacement := lipgloss.Place(m.width-52, m.height, lipgloss.Center, lipgloss.Right, gifBox)
  	banmenu := lipgloss.JoinVertical(lipgloss.Top, bannerStyle.Render(banner), menuPlacement)
  	totalscreen := lipgloss.JoinHorizontal(lipgloss.Left, banmenu, gifPlacement)

  return totalscreen 
}