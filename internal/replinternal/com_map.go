package replinternal

import (
	"fmt"
	"errors"
	"strings"
	"github.com/som-pat/poke_dex/internal/config"
)

func call_map(cfg_state *config.ConfigState, args ...string) (string,[]string,error){
	resp, err := cfg_state.PokeapiClient.InvokeLocs(cfg_state.NextLocURL)
	if err!= nil{
		return "",nil,errors.New("next page does not exist")
	}
	lislocs := make([]string, 0, 20)
	var displayloc strings.Builder
	displayloc.WriteString("Regions: \n")
	for _, area := range resp. Results{
		lislocs = append(lislocs,area.Name)
		displayloc.WriteString(fmt.Sprintf("- %s\n", area.Name))
	}
	cfg_state.NextLocURL = resp.Next
	cfg_state.PrevLocURL = resp.Previous
	return displayloc.String(),lislocs,nil
}



func call_mapb(cfg_state *config.ConfigState, args ...string) (string,[]string,error){	
	if cfg_state.PrevLocURL == nil{
		return "",nil,errors.New("you're on the 1st page")
	}
	resp, err := cfg_state.PokeapiClient.InvokeLocs(cfg_state.PrevLocURL)
	if err!= nil{
		return "",nil,err
	}
	lislocs := make([]string, 0, 20)
	var displayprevloc strings.Builder
	displayprevloc.WriteString("Previous Regions: \n")
	for _, area := range resp. Results{
		lislocs = append(lislocs,area.Name)
		displayprevloc.WriteString(fmt.Sprintf("- %s\n", area.Name))
	}
	cfg_state.NextLocURL = resp.Next
	cfg_state.PrevLocURL = resp.Previous

	return displayprevloc.String(),lislocs,nil
}