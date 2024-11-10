package main

import (
	"errors"
	"fmt"
	"strings"
	"bytes"
	"os/exec"
)


func call_pokeInspect(cfg_state *ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no pokemon or item name provided")
	}	
	toInspect := args[0]	

	pokemon, ok := cfg_state.pokemonCaught[toInspect] 
	if ok {
		var pokedetails strings.Builder
		ascii_img,err := getPokemonAscii(pokemon.Sprites.FrontDefault)
		if err!= nil{
			pokedetails.WriteString(" [Image Unavailable]\n")
		}
		pokedetails.WriteString(ascii_img +"\n")

		pokedetails.WriteString(fmt.Sprintf("Name: %s \n", pokemon.Name))
		pokedetails.WriteString(fmt.Sprintf("Height: %d \n", pokemon.Height))
		pokedetails.WriteString(fmt.Sprintf("Weight:%d \n", pokemon.Weight))
		pokedetails.WriteString("Stats:\n")
		for _, stat := range pokemon.Stats {
			pokedetails.WriteString(fmt.Sprintf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat))
		}
		pokedetails.WriteString("Types:")
		for _, typeInfo := range pokemon.Types {
			pokedetails.WriteString(fmt.Sprintf(" %s,", typeInfo.Type.Name))
		}
		
		return pokedetails.String(),nil,nil
	}
	helditem, ok := cfg_state.ItemsHeld[toInspect]
	if !ok{
		return "",nil,fmt.Errorf("%s not in inventory, can't inspect", toInspect)
	}
	var itemdetails strings.Builder
	ascii_img,err := getPokemonAscii(helditem.Sprites.Default)
	if err!= nil{
		itemdetails.WriteString(" [Image Unavailable]\n")
	}
	itemdetails.WriteString(ascii_img +"\n")
	itemdetails.WriteString(fmt.Sprintf("Name: %s \n", helditem.Name))
	itemdetails.WriteString(fmt.Sprintf("Category: %s \n", helditem.Category.Name))
	itemdetails.WriteString(fmt.Sprintf("Cost: %d \n",helditem.Cost))
	itemdetails.WriteString(fmt.Sprintf("Effect : %s \n",helditem.EffectEntries[0].ShortEffect))
	itemdetails.WriteString("Attributes: \n")
	for _, attr := range helditem.Attributes{
		itemdetails.WriteString(fmt.Sprintf(" -%s \n",attr.Name))
	}
	return itemdetails.String(), nil, nil
}

func getPokemonAscii(imageURL string) (string, error){
	cmd := exec.Command("./ascii/bin/python3","ascii_py/new_ascii.py", imageURL)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err!=nil{
		return "", err
	}

	asciiImg := out.String()
	return asciiImg, nil
}

