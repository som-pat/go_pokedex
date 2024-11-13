package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) InvokeMove(moveurl string) (FilPokeMove, error){
	full_url := moveurl
	
	req, err := http.NewRequest("GET", full_url, nil)
	if err!= nil {return FilPokeMove{}, err}

	resp,err := c.httpClient.Do(req)
	if err!= nil {return FilPokeMove{}, err}

	defer resp.Body.Close()

	if resp.StatusCode > 399{
		return FilPokeMove{}, fmt.Errorf("bad Status Encounterd: %v", resp.StatusCode)
	}
	data, err:= io.ReadAll(resp.Body)
	if err!=nil{
		return FilPokeMove{},err
	}

	pokemove := FilPokeMove{}
	err = json.Unmarshal(data, &pokemove)
	if err != nil{
		return FilPokeMove{}, err
	}
	return pokemove, nil


}