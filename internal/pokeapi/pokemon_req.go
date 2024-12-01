package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func (c *Client) InvokePokeCatch(Poke_name interface{}) (PokemonDetails,error){
	var nameid string
	switch v := Poke_name.(type){
	case string:
		nameid = v
	case int:
		nameid = fmt.Sprintf("%d",v)
	default:
		return PokemonDetails{}, fmt.Errorf("not proper type")
	}
	end_point := "/pokemon/" + nameid
	full_url := baseURL + end_point


	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return PokemonDetails{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonDetails{}, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return PokemonDetails{}, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDetails{}, err
	}
	
	poke_details := PokemonDetails{}
	err = json.Unmarshal(data, &poke_details)
	
	if err != nil {
		return PokemonDetails{}, err
	}

	return poke_details, nil
}



func (c *Client) EncounterPoke(Poke_name string) (FilPokemonSpecies,error){
	end_point := "/pokemon-species/" + Poke_name
	full_url := baseURL + end_point
	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return FilPokemonSpecies{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return FilPokemonSpecies{}, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return FilPokemonSpecies{}, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return FilPokemonSpecies{}, err
	}
	
	poke_species_details := FilPokemonSpecies{}
	err = json.Unmarshal(data, &poke_species_details)
	
	if err != nil {
		return FilPokemonSpecies{}, err
	}
	return poke_species_details, nil
}

