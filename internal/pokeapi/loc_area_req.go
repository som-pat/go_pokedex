package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func (c *Client) InvokeLocs(pageURL *string) (LocationAreaResp, error){
	end_point := "/location-area"
	full_url := baseURL + end_point
	if pageURL != nil {
		full_url = *pageURL
	}
	// check cache
	cache_data, ok := c.cache.Get(full_url)
	if ok{
		loc_resp := LocationAreaResp{}
		err := json.Unmarshal(cache_data, &loc_resp)		
		if err != nil {
			return LocationAreaResp{}, err
		}
		return loc_resp, nil
	}

	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return LocationAreaResp{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreaResp{}, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return LocationAreaResp{}, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}
	
	loc_resp := LocationAreaResp{}
	err = json.Unmarshal(data, &loc_resp)
	
	if err != nil {
		return LocationAreaResp{}, err
	}

	c.cache.Add(full_url, data)
	return loc_resp, nil


}



func (c *Client) InvokePokeLocs(LocatioName string) (PokeinLoc, error){
	end_point := "/location-area/" + LocatioName
	full_url := baseURL + end_point
	// check cache
	cache_data, ok := c.cache.Get(full_url)
	if ok{
		fmt.Println("Cache hit, looting booty")
		poke_loc_resp := PokeinLoc{}
		err := json.Unmarshal(cache_data, &poke_loc_resp)		
		if err != nil {
			return PokeinLoc{}, err
		}
		return poke_loc_resp, nil
	}

	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return PokeinLoc{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokeinLoc{}, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return PokeinLoc{}, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeinLoc{}, err
	}
	
	poke_loc_resp := PokeinLoc{}
	err = json.Unmarshal(data, &poke_loc_resp)
	
	if err != nil {
		return PokeinLoc{}, err
	}

	c.cache.Add(full_url, data)
	return poke_loc_resp, nil


}
