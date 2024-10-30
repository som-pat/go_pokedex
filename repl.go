package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config_state, ...string) error
}

func get_command() map[string] cliCommand{
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    call_help,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    call_exit,
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
			name:		 "Inspect{Pokemon_name}",
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

func repl_input(cfg_state *config_state){
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		input.Scan()
		input_text := input.Text()
		new_input := input_clean(input_text)
		
		// Empty commands
		if len(new_input) == 0{
			continue
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
			fmt.Println("Not a generic command")
			continue
		} 
		err :=route_com.callback(cfg_state, args...)
		if err != nil {
			fmt.Println(err)
		}

	}	
}

func input_clean(input string) ([]string) {
	lower :=  strings.ToLower(input)
	words :=  strings.Fields(lower)
	return words
}