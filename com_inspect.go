package main

import (
	"errors"
	"fmt"
	"strings"
)


func call_pokeInspect(cfg_state *ConfigState, args ...string) (string,error){
	if len(args) != 1{
		return "",errors.New("no pokemon name provided")
	}	
	poke_name := args[0]	

	pokemon, ok := cfg_state.pokemonCaught[poke_name] 
	if !ok {
		return "",fmt.Errorf("pokemon %s not caught", poke_name)
	}
	var pokedetails strings.Builder
	pokedetails.WriteString(fmt.Sprintf("Name: %s \n", pokemon.Name))
	pokedetails.WriteString(fmt.Sprintf("Height: %d \n", pokemon.Height))
	pokedetails.WriteString(fmt.Sprintf("Weight:%d \n", pokemon.Weight))
	pokedetails.WriteString("Stats:\n")
	for _, stat := range pokemon.Stats {
		pokedetails.WriteString(fmt.Sprintf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat))
	}
	pokedetails.WriteString("Types:")
	for _, typeInfo := range pokemon.Types {
		pokedetails.WriteString(fmt.Sprintf(" %s,", typeInfo.Type.Name))
	}

	
	return pokedetails.String(), nil
}