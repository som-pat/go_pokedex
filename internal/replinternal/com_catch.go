package replinternal

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"github.com/som-pat/poke_dex/internal/config"
)

func call_catch(cfg_state *config.ConfigState, args ...string) (string,[]string,error){
	if len(args) != 1{
		return "",nil,errors.New("no pokemon or item name provided")
	}
	toScour := args[0]
	_,exists := cfg_state.PokemonCaught[toScour]
	if exists{
		return "",nil, fmt.Errorf("%s already exists in your inventory",toScour)
	}
	pokemon, err := cfg_state.PokeapiClient.InvokePokeCatch(toScour)
	if err != nil{
		return "", nil, errors.New("no such pokemon found")
	}else{	
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
				cfg_state.PokemonCaught[pokemon.Name] = pokemon
				break
			}
		}
		pokemon, ok := cfg_state.PokemonCaught[pokemon.Name]
		if !ok{
			catchchance.WriteString("\n")
			catchchance.WriteString(fmt.Sprintf("Unable to catch %s, better luck next time \n", pokemon.Name))
			catchchance.WriteString("\n")
		}
		return catchchance.String(),nil,nil
	}
}


