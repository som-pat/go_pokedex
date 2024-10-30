package main

import "fmt"
import "errors"


func call_pokedex(cfg_state *config_state, args ...string) error{
	if len(cfg_state.pokemonCaught) == 0{
		return errors.New("no pokemon caught till now")
	}

	for _, poke := range cfg_state.pokemonCaught{
		fmt.Printf(" - %s \n", poke.Name)
	}
	
	return nil
}