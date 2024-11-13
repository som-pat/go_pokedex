package main

import (
	"fmt"
	"log"
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
	f, err:= tea.LogToFile("debug.log", "debug")
	if err != nil{ log.Fatalf("Error encountered %v",err)}
	defer f.Close()

    p:= tea.NewProgram(takeInput(cfgState),tea.WithAltScreen())
    if _,err := p.Run(); err != nil {
        fmt.Printf("Error starting program: %v\n", err)
    }
}


