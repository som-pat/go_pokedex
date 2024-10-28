package main

import (
	"fmt"
	"log"
)

func call_mapb(cfg_state *config_state) error{	
	resp, err := cfg_state.pokeapiClient.InvokeLocs(cfg_state.prevLocURL)
	if err!= nil{
		log.Fatal(err)
	}
	fmt.Println("Previous Locations:")
	for _, area := range resp. Results{
		fmt.Printf("- %s\n", area.Name)
	}
	cfg_state.nextLocURL = resp.Next
	cfg_state.prevLocURL = resp.Previous

	return nil
}