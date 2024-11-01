package main

import (
	"fmt"
	"errors"
	"strings"
)
func call_map(cfg_state *ConfigState, args ...string) (string,[]string,error){
	resp, err := cfg_state.pokeapiClient.InvokeLocs(cfg_state.nextLocURL)
	if err!= nil{
		return "",nil,errors.New("next page does not exist")
	}
	var displayloc strings.Builder
	displayloc.WriteString("Locations: \n")
	for _, area := range resp. Results{
		displayloc.WriteString(fmt.Sprintf("- %s\n", area.Name))
	}
	cfg_state.nextLocURL = resp.Next
	cfg_state.prevLocURL = resp.Previous
	return displayloc.String(),nil,nil
}



func call_mapb(cfg_state *ConfigState, args ...string) (string,[]string,error){	
	if cfg_state.prevLocURL == nil{
		return "",nil,errors.New("you're on the 1st page")
	}
	resp, err := cfg_state.pokeapiClient.InvokeLocs(cfg_state.prevLocURL)
	if err!= nil{
		return "",nil,err
	}
	var displayprevloc strings.Builder
	displayprevloc.WriteString("Previous Locations: \n")
	for _, area := range resp. Results{
		displayprevloc.WriteString(fmt.Sprintf("- %s\n", area.Name))
	}
	cfg_state.nextLocURL = resp.Next
	cfg_state.prevLocURL = resp.Previous

	return displayprevloc.String(),nil,nil
}