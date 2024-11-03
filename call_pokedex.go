package main

import "fmt"
import "errors"
import "strings"


func call_pokedex(cfg_state *ConfigState, args ...string) (string,[]string,error){
	if len(cfg_state.pokemonCaught) == 0{
		return "",nil,errors.New("no pokemon caught till now")
	}
	var caught_pokemon strings.Builder
	var cp []string
	for _, poke := range cfg_state.pokemonCaught{
		caught_pokemon.WriteString(fmt.Sprintf(" - %s \n", poke.Name))
		cp = append(cp, poke.Name)
	}
	
	return caught_pokemon.String(),cp,nil
}