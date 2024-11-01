package main

import (
	tea "github.com/charmbracelet/bubbletea"

)

func Run(cfgState *ConfigState) {
    p:= tea.NewProgram(takeInput(cfgState))
    if _,err := p.Run(); err != nil {
        panic(err)
    }
}


