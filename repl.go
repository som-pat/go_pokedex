package main

import (
	"fmt"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*ConfigState, ...string) (string,[]string,error)
}

func get_command() map[string] cliCommand{
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays available commands",
			callback:    call_help,
		},

		"map":{
			name:		 "map",
			description: "Display next 20 loactions",
			callback:	 call_map,
		},
		
		"mapb":{
			name:		 "mapb",
			description: "Display previous 20 locations",
			callback:    call_mapb,
		},
		
		"explore":{
			name:		 "explore 'location_area' ",
			description: "Display Pokemons in chosen location",
			callback:    call_explore,
		},
		"scout":{
			name:		 "scout 'pokemon/item' ",
			description: "Search for given Pokemons/items in chosen location",
			callback:    call_search,
		},

		"catch":{
			name:		 "catch{Pokemon/Item name}",
			description: "Catch Pokemons",
			callback:    call_catch,
		},
		
		"inspect":{
			name:		 "inspect{Pokemon/Item name}",
			description: "Inspect caught Pokemons",
			callback:    call_pokeInspect,
		},

		"inventory":{
			name:		 "inventory",
			description: "View caught Pokemons/ held Items",
			callback:    call_inventory,
		},
	}
}

func repl_input(cfg_state *ConfigState, input string) (string,[]string){
	new_input := input_clean(input)
	
	// Empty commands
	if len(new_input) == 0{
		return "No command entered. Type 'help' for a list of commands.",nil
	}
	com := new_input[0]
	args := []string{}
	if len(new_input)>1{
		args = new_input[1:]
	}

	avail_com := get_command()
	
	// Check if valid command
	route_com,ok  := avail_com[com]
	if !ok{
		return "Unknown command. Type 'help' for a list of available commands.",nil
	}

	res, lis, err :=route_com.callback(cfg_state, args...)
	if err != nil {
		return fmt.Sprintf("Error: %v",err),nil
	}

	return res,lis
		
}

func input_clean(input string) ([]string) {
	lower :=  strings.ToLower(input)
	words :=  strings.Fields(lower)
	return words
}