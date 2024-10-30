package main

import (
	"time"

	"github.com/som-pat/poke_dex/internal/pokeapi"
)

type config_state struct{
	pokeapiClient pokeapi.Client
	nextLocURL *string
	prevLocURL *string
	pokemonCaught map[string] pokeapi.PokemonDetails
}

func main()	{
	cfg_state := config_state{
		pokeapiClient: pokeapi.NewClient(time.Hour),
		pokemonCaught: make(map[string]pokeapi.PokemonDetails),
	}	
	repl_input(&cfg_state)
	
}