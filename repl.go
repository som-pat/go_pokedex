package main

import (
	"fmt"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*ConfigState, ...string) (string,error)
}

func get_command() map[string] cliCommand{
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays available commands",
			callback:    call_help,
		},
		// "exit": {
		// 	name:        "exit",
		// 	description: "Exit the Pokedex",
		// 	callback:    call_exit,
		// },
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
			name:		 "explore{Location_area}",
			description: "Display Pokemons in chosen location",
			callback:    call_explore,
		},

		"catch":{
			name:		 "catch{Pokemon_name}",
			description: "Catch Pokemons",
			callback:    call_catch,
		},
		
		"inspect":{
			name:		 "inspect{Pokemon_name}",
			description: "Inspect caught Pokemons",
			callback:    call_pokeInspect,
		},

		"pokedex":{
			name:		 "pokedex",
			description: "View caught Pokemons",
			callback:    call_pokedex,
		},
	}
}

func repl_input(cfg_state *ConfigState, input string) (string){
	new_input := input_clean(input)
	
	// Empty commands
	if len(new_input) == 0{
		return "No command entered. Type 'help' for a list of commands."
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
		return "Unknown command. Type 'help' for a list of available commands."
	}

	res,err :=route_com.callback(cfg_state, args...)
	if err != nil {
		return fmt.Sprintf("Error: %v",err)
	}

	return res
		
}

func input_clean(input string) ([]string) {
	lower :=  strings.ToLower(input)
	words :=  strings.Fields(lower)
	return words
}