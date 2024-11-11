package config

import (
	"github.com/som-pat/poke_dex/internal/pokeapi"
)


type ConfigState struct{
	PokeapiClient 		 pokeapi.Client
	NextLocURL 	  		 *string
	PrevLocURL	  		 *string
	PokemonCaught 		 map[string] pokeapi.PokemonDetails
	ItemsHeld 	  		 map[string] pokeapi.ItemDescription
	CurrentEncounterList *[]string
}