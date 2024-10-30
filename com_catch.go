package main

import "fmt"
import "errors"
import "math/rand"

func call_catch(cfg_state *config_state, args ...string) error{
	if len(args) != 1{
		return errors.New("no pokemon name provided")
	}
	poke_name := args[0]	

	pokemon, err := cfg_state.pokeapiClient.InvokePokeCatch(poke_name)
	if err!= nil{
		return err
	}
	
	// three chances to catch after that poke will escape
	const chances = 60
	for i:=1;i<=3;i++{
		randChances := rand.Intn(chances)
		randBaseExp := rand.Intn(pokemon.BaseExperience)
		if randChances < randBaseExp{
			fmt.Printf("%s not caught \n", pokemon.Name)
		}else{
		fmt.Printf("%s caught \n", pokemon.Name)
		cfg_state.pokemonCaught[pokemon.Name] = pokemon
		break
		}
	fmt.Println()
	fmt.Printf("Unable to catch %s, better luck next time \n", poke_name)
	} 
	return nil
}