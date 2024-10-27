package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func (c *Client) InvokeLocs() (LocationAreaResp, error){
	end_point := "/location"
	full_url := baseURL + end_point

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
	return loc_resp, nil


}