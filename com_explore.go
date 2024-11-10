package main

import (
	"errors"
	"fmt"
	"strings"
	"math/rand"
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

	const NumItems = 8
	var Itemsadd = rand.Intn(NumItems)
	itemName, err := cfg_state.pokeapiClient.ItemRandomizer(Itemsadd)
	if err != nil {
		return "",nil,err
	}

	var explore_reg strings.Builder
	var lisexp_reg []string
	explore_reg.WriteString(fmt.Sprintf("Items and Pokemons found in %s:\n", location.Name))
	for _, poke_struct := range location.PokemonEncounters {
		lisexp_reg = append(lisexp_reg, poke_struct.Pokemon.Name)
		explore_reg.WriteString(fmt.Sprintf("- %s\n", poke_struct.Pokemon.Name))
	}
	

	for _, iname := range(itemName){
		lisexp_reg = append(lisexp_reg, iname)
		explore_reg.WriteString(fmt.Sprintf("- %s\n",iname))

	}

	if len(lisexp_reg) == 0{
		explore_reg.WriteString(fmt.Sprintf("No Pokemons or Items found in %s:\n", location.Name))
		return explore_reg.String(),nil,nil
	}
	return explore_reg.String(),lisexp_reg,nil
}
