package main

import "fmt"
import "errors"


func call_pokeInspect(cfg_state *config_state, args ...string) error{
	if len(args) != 1{
		return errors.New("no pokemon name provided")
	}	
	poke_name := args[0]	

	pokemon, ok := cfg_state.pokemonCaught[poke_name] 
	if !ok {
		return fmt.Errorf("pokemon %s not caught", poke_name)
	}
	
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Println("  -", typeInfo.Type.Name)
	}

	
	return nil
}