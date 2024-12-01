package config

import "github.com/som-pat/poke_dex/internal/pokeapi"

type userbattle struct{
	PlayerName			string
	PokemonRoster		[]string
	PokeLvl				[]int
	RosterDetails		map[string] pokeapi.FilPokemonSpecies // pokemonDetails check
}

