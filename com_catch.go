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
	
	const threshold = 50
	
		randchances := rand.Intn(pokemon.BaseExperience)
		if randchances > threshold{
			return fmt.Errorf("%s not caught ", pokemon.Name)
		}
		fmt.Printf("%s caught \n", pokemon.Name) 
	
	return nil
}