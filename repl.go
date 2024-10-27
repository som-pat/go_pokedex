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
	callback    func() error
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
	}
}

func repl_input(){
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex >")
		input.Scan()
		input_text := input.Text()
		new_input := input_clean(input_text)
		if len(new_input) == 0{
			continue
		}
		com := new_input[0]
		avail_com := get_command()
		
		route_com,ok  := avail_com[com]
		if !ok{
			fmt.Println("Not a generic command")
			continue
		} 
		route_com.callback()

	}	
}

func input_clean(input string) ([]string) {
	lower :=  strings.ToLower(input)
	words :=  strings.Fields(lower)
	return words
}