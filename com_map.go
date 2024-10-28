package main

import (
	"fmt"
	"errors"
)
func call_map(cfg_state *config_state) error{
	resp, err := cfg_state.pokeapiClient.InvokeLocs(cfg_state.nextLocURL)
	if err!= nil{
		return err
	}
	fmt.Println("Locations:")
	for _, area := range resp. Results{
		fmt.Printf("- %s\n", area.Name)
	}
	cfg_state.nextLocURL = resp.Next
	cfg_state.prevLocURL = resp.Previous
	return nil
}

func call_mapb(cfg_state *config_state) error{	
	if cfg_state.prevLocURL == nil{
		return errors.New("you're on the 1st page")
	}
	resp, err := cfg_state.pokeapiClient.InvokeLocs(cfg_state.prevLocURL)
	if err!= nil{
		return err
	}
	fmt.Println("Previous Locations:")
	for _, area := range resp. Results{
		fmt.Printf("- %s\n", area.Name)
	}
	cfg_state.nextLocURL = resp.Next
	cfg_state.prevLocURL = resp.Previous

	return nil
}