package main

import (
	"errors"
	"fmt"
	"strings"
)

func call_explore(cfg_state *ConfigState, args ...string) (string,[]string,error) {
	if len(args) != 1 {
		return "",nil,errors.New("no Region provided")
	}
	loc_name := args[0]

	location, err := cfg_state.pokeapiClient.InvokePokeLocs(loc_name)
	if err != nil {
		return "",nil,err
	}
	var explore_reg strings.Builder
	var lisexp_reg []string
	explore_reg.WriteString(fmt.Sprintf("Pokemons in %s:\n", location.Name))
	for _, poke_struct := range location.PokemonEncounters {
		lisexp_reg = append(lisexp_reg, poke_struct.Pokemon.Name)
		explore_reg.WriteString(fmt.Sprintf("- %s\n", poke_struct.Pokemon.Name))
	}
	if len(lisexp_reg) == 0{
		explore_reg.WriteString(fmt.Sprintf("No Pokemons found in %s:\n", location.Name))
		return explore_reg.String(),nil,nil
	}
	return explore_reg.String(),lisexp_reg,nil
}
