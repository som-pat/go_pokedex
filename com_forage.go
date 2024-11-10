package main

import (
	"errors"
	"fmt"
	"strings"
)


func call_forage(cfg_state *ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no item name provided")
	}
	toForage := args[0]
	itemdes, itemErr := cfg_state.pokeapiClient.ItemFetch(toForage)
	if itemErr != nil{
		return "", nil, errors.New("no such item found")
	}
	var itemScoured strings.Builder
	cfg_state.ItemsHeld[itemdes.Name] = itemdes
	itemScoured.WriteString(fmt.Sprintf("Foraged %s \n",itemdes.Name))
	itemScoured.WriteString(fmt.Sprintf("len of %d \n",len(cfg_state.ItemsHeld)))
	return itemScoured.String(),nil, nil	

}