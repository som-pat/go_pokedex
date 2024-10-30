package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func (c *Client) InvokePokeCatch(Poke_name string) (PokemonDetails,error){
	end_point := "/pokemon/" + Poke_name
	full_url := baseURL + end_point

	// check cache
	cache_data, ok := c.cache.Get(full_url)
	if ok{
		fmt.Println("Cache hit, looting booty")
		poke_details := PokemonDetails{}
		err := json.Unmarshal(cache_data, &poke_details)		
		if err != nil {
			return PokemonDetails{}, err
		}
		return poke_details, nil
	}

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

	c.cache.Add(full_url, data)
	return poke_details, nil


}



