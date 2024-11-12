package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/som-pat/poke_dex/internal/config"
	"github.com/som-pat/poke_dex/internal/pokeapi"
)



func main()	{
	
	cfg_state := &config.ConfigState{
		PokeapiClient: pokeapi.NewClient(time.Hour),
		PokemonCaught: make(map[string]pokeapi.PokemonDetails),
		ItemsHeld: make(map[string]pokeapi.ItemDescription),
		CurrentEncounterList: &[]string{},
	}

	Run(cfg_state)
	
}

func Run(cfgState *config.ConfigState) {
    p:= tea.NewProgram(takeInput(cfgState))
    if _,err := p.Run(); err != nil {
        fmt.Printf("Error starting program: %v\n", err)
    }
}


