package main

import (
	"errors"
	"fmt"
	"strings"
)

func call_inventory(cfg_state *ConfigState, args ...string) (string,[]string,error){
	if len(cfg_state.pokemonCaught) ==0 && len(cfg_state.ItemsHeld)==0{
		return "",nil,errors.New("no pokemons or items in the inventory")
	}
	var held_poke_item strings.Builder
	var cp []string
	if len(cfg_state.pokemonCaught) >0 {
		held_poke_item.WriteString("Pokemons Held:\n")
		for _, poke := range cfg_state.pokemonCaught{
			held_poke_item.WriteString(fmt.Sprintf(" - %s \n", poke.Name))
			cp = append(cp, poke.Name)		
		}
	}
	if len(cfg_state.ItemsHeld)>0{
		held_poke_item.WriteString("Items Held:\n")
		for _, item := range cfg_state.ItemsHeld{
				held_poke_item.WriteString(fmt.Sprintf(" - %s \n", item.Name))
				cp = append(cp, item.Name)
			}
	}
	
	return held_poke_item.String(),cp,nil
}


