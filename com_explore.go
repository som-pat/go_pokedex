package main

import (
	"errors"
	"fmt"
	"strings"
)

func call_explore(cfg_state *ConfigState, args ...string) (string,error) {
	if len(args) != 1 {
		return "",errors.New("no location area provided")
	}
	loc_name := args[0]

	location, err := cfg_state.pokeapiClient.InvokePokeLocs(loc_name)
	if err != nil {
		return "",err
	}
	var explore_reg strings.Builder
	explore_reg.WriteString(fmt.Sprintf("Pokemons in %s:\n", location.Name))
	for _, poke_struct := range location.PokemonEncounters {
		explore_reg.WriteString(fmt.Sprintf("- %s\n", poke_struct.Pokemon.Name))

	}
	return explore_reg.String(),nil
}
