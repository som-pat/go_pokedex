
package app

import tea "github.com/charmbracelet/bubbletea"


type Navigator interface {
	GoToMenu() tea.Model
	GoToStoryMode() tea.Model
	GoToBattleMode() tea.Model
}

type AppNavigator struct {
	menu    tea.Model
	story   tea.Model
	battle  tea.Model
}

func NewAppNavigator(menu, story, battle tea.Model) *AppNavigator {
	return &AppNavigator{
		menu:    menu,
		story:   story,
		battle:  battle,
	}
}


func (n *AppNavigator) GoToMenu() tea.Model {
	return n.menu
}

func (n *AppNavigator) GoToStoryMode() tea.Model {
	return n.story
}

func (n *AppNavigator) GoToBattleMode() tea.Model {
	return n.battle
}
