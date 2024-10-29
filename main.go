package main

import (
	"time"

	"github.com/som-pat/poke_dex/internal/pokeapi"
)

type config_state struct{
	pokeapiClient pokeapi.Client
	nextLocURL *string
	prevLocURL *string
}

func main()	{
	cfg_state := config_state{
		pokeapiClient: pokeapi.NewClient(time.Hour),
	}	
	repl_input(&cfg_state)
	
}