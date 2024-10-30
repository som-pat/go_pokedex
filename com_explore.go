package main

import (
	"errors"
	"fmt"
)

func call_explore(cfg_state *config_state, args ...string) error {
	if len(args) != 1 {
		return errors.New("no location area provided")
	}
	loc_name := args[0]

	location, err := cfg_state.pokeapiClient.InvokePokeLocs(loc_name)
	if err != nil {
		return err
	}
	fmt.Printf("Pokemons in %s:\n", location.Name)
	for _, poke_struct := range location.PokemonEncounters {
		fmt.Printf("- %s\n", poke_struct.Pokemon.Name)

	}
	return nil
}
