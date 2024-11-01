package main

import (
	"time"
	"github.com/som-pat/poke_dex/internal/pokeapi"


)

type ConfigState struct{
	pokeapiClient pokeapi.Client
	nextLocURL *string
	prevLocURL *string
	pokemonCaught map[string] pokeapi.PokemonDetails
}


func main()	{

	cfg_state := ConfigState{
		pokeapiClient: pokeapi.NewClient(time.Hour),
		pokemonCaught: make(map[string]pokeapi.PokemonDetails),
	}

	Run(&cfg_state)
		
	// repl_input(&cfg_state)
	
}