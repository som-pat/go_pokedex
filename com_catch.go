package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

func call_catch(cfg_state *ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no pokemon name provided")
	}
	poke_name := args[0]	

	pokemon, err := cfg_state.pokeapiClient.InvokePokeCatch(poke_name)
	if err!= nil{
		return "",nil,err
	}
	
	var catchchance strings.Builder
	// three chances to catch after that poke will escape
	const chances = 60
	for i:=1;i<=3;i++{
		randChances := rand.Intn(chances)
		randBaseExp := rand.Intn(pokemon.BaseExperience)
		if randChances < randBaseExp{
			catchchance.WriteString(fmt.Sprintf("%s not caught \n", pokemon.Name))
		}else{
		catchchance.WriteString(fmt.Sprintf("%s caught \n", pokemon.Name))
		cfg_state.pokemonCaught[pokemon.Name] = pokemon
		break
		}
	}
	pokemon, ok := cfg_state.pokemonCaught[pokemon.Name]
	if !ok{
		catchchance.WriteString("\n")
		catchchance.WriteString(fmt.Sprintf("Unable to catch %s, better luck next time \n", poke_name))
		catchchance.WriteString("\n")
	}
	return catchchance.String(),nil,nil
}